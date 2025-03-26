package useCase

import (
	"log/slog"
	"regexp"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

const SslValidationsPerHour int = 4

type SslCertificateWatchdog struct {
	vHostQueryRepo repository.VirtualHostQueryRepo
	sslQueryRepo   repository.SslQueryRepo
	sslCmdRepo     repository.SslCmdRepo
}

func NewSslCertificateWatchdog(
	vHostQueryRepo repository.VirtualHostQueryRepo,
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
) *SslCertificateWatchdog {
	return &SslCertificateWatchdog{
		vHostQueryRepo: vHostQueryRepo,
		sslQueryRepo:   sslQueryRepo,
		sslCmdRepo:     sslCmdRepo,
	}
}

type sslWatchdogVirtualHostSettings struct {
	VirtualHost    entity.VirtualHost
	AliasHostnames []valueObject.Fqdn
	SslPair        entity.SslPair
}

func (uc *SslCertificateWatchdog) vHostSettingsMapFactory() (
	vhostSettingsMap map[valueObject.Fqdn]*sslWatchdogVirtualHostSettings,
) {
	vhostEntities, err := uc.vHostQueryRepo.Read()
	if err != nil {
		slog.Error(
			"ReadVirtualHostInfraError",
			slog.String("error", err.Error()),
			slog.String("method", "SslCertificateWatchdog"),
		)
		return vhostSettingsMap
	}

	for _, vhostEntity := range vhostEntities {
		mainHostname := vhostEntity.Hostname
		isAlias := vhostEntity.Type == valueObject.VirtualHostTypeAlias
		if isAlias {
			if vhostEntity.ParentHostname == nil {
				slog.Debug(
					"AliasWithoutParentHostname",
					slog.String("alias", vhostEntity.Hostname.String()),
					slog.String("method", "SslCertificateWatchdog"),
				)
				continue
			}

			mainHostname = *vhostEntity.ParentHostname
		}

		vhostMap, mapExists := vhostSettingsMap[mainHostname]
		if !mapExists {
			vhostMap = &sslWatchdogVirtualHostSettings{}
		}

		if isAlias {
			vhostMap.AliasHostnames = append(vhostMap.AliasHostnames, vhostEntity.Hostname)
			continue
		}

		vhostMap.VirtualHost = vhostEntity
	}

	sslPairEntities, err := uc.sslQueryRepo.Read()
	if err != nil {
		slog.Error(
			"ReadSslPairInfraError",
			slog.String("error", err.Error()),
			slog.String("method", "SslCertificateWatchdog"),
		)
		return vhostSettingsMap
	}

	for _, sslPairEntity := range sslPairEntities {
		vhostMap, mapExists := vhostSettingsMap[sslPairEntity.MainVirtualHostHostname]
		if !mapExists {
			slog.Debug(
				"SslPairWithUnknownVirtualHost",
				slog.String("sslPairId", sslPairEntity.Id.String()),
				slog.String("method", "SslCertificateWatchdog"),
			)
			continue
		}

		vhostMap.SslPair = sslPairEntity
	}

	return vhostSettingsMap
}

func (uc *SslCertificateWatchdog) Execute() {
	vhostSettingsMap := uc.vHostSettingsMapFactory()

	skipSubdomainRegex := regexp.MustCompile(`^[^\.]+\.`)
	for _, vhostSettings := range vhostSettingsMap {
		shouldRenewCert := vhostSettings.SslPair.IsPubliclyTrusted()

		certAltNamesStrMap := map[string]interface{}{}
		for _, altName := range vhostSettings.SslPair.Certificate.AltNames {
			certAltNamesStrMap[altName.String()] = nil
		}

		missingAltNames := []valueObject.Fqdn{}
		for _, aliasHostname := range vhostSettings.AliasHostnames {
			aliasHostnameStr := aliasHostname.String()
			_, altNameExists := certAltNamesStrMap[aliasHostname.String()]
			if altNameExists {
				continue
			}

			likelyParentHostname := skipSubdomainRegex.ReplaceAllString(aliasHostnameStr, "")
			wildcardParentHostname := "*." + likelyParentHostname
			_, altNameExists = certAltNamesStrMap[wildcardParentHostname]
			if altNameExists {
				continue
			}

			missingAltNames = append(missingAltNames, aliasHostname)
		}

		if len(missingAltNames) > 0 {
			shouldRenewCert = true
		}

		if !shouldRenewCert {
			continue
		}

		err := uc.sslCmdRepo.CreatePubliclyTrusted(dto.NewCreatePubliclyTrustedSslPair(
			vhostSettings.VirtualHost.Hostname, vhostSettings.AliasHostnames,
			valueObject.AccountIdSystem, valueObject.IpAddressSystem,
		))
		if err != nil {
			slog.Error(
				"CreatePubliclyTrustedSslPairError",
				slog.String("error", err.Error()),
				slog.String("hostname", vhostSettings.VirtualHost.Hostname.String()),
				slog.String("sslPairId", vhostSettings.SslPair.Id.String()),
				slog.String("method", "SslCertificateWatch"),
			)
		}
	}
}
