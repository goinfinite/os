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
	vhostQueryRepo repository.VirtualHostQueryRepo
	sslQueryRepo   repository.SslQueryRepo
	sslCmdRepo     repository.SslCmdRepo
}

func NewSslCertificateWatchdog(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	sslQueryRepo repository.SslQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
) *SslCertificateWatchdog {
	return &SslCertificateWatchdog{
		vhostQueryRepo: vhostQueryRepo,
		sslQueryRepo:   sslQueryRepo,
		sslCmdRepo:     sslCmdRepo,
	}
}

type sslWatchdogVirtualHostSettings struct {
	VirtualHost    entity.VirtualHost
	AliasHostnames []valueObject.Fqdn
	SslPair        entity.SslPair
}

func (uc *SslCertificateWatchdog) vHostSettingsMapFactory() map[valueObject.Fqdn]*sslWatchdogVirtualHostSettings {
	vhostSettingsMap := map[valueObject.Fqdn]*sslWatchdogVirtualHostSettings{}

	readVirtualHostsResponse, err := uc.vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
		Pagination: dto.PaginationUnpaginated,
	})
	if err != nil {
		slog.Error(
			"ReadVirtualHostInfraError",
			slog.String("error", err.Error()),
			slog.String("method", "SslCertificateWatchdog"),
		)
		return vhostSettingsMap
	}

	virtualHostAliases := []entity.VirtualHost{}
	for _, vhostEntity := range readVirtualHostsResponse.VirtualHosts {
		if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
			virtualHostAliases = append(virtualHostAliases, vhostEntity)
			continue
		}

		vhostSettingsMap[vhostEntity.Hostname] = &sslWatchdogVirtualHostSettings{
			VirtualHost:    vhostEntity,
			AliasHostnames: []valueObject.Fqdn{},
			SslPair:        entity.SslPair{},
		}
	}

	for _, vhostAliasEntity := range virtualHostAliases {
		if vhostAliasEntity.ParentHostname == nil {
			slog.Debug(
				"AliasVirtualHostWithoutParent",
				slog.String("method", "SslCertificateWatchdog"),
				slog.String("aliasHostname", vhostAliasEntity.Hostname.String()),
			)
			continue
		}

		vhostSettings, mapExists := vhostSettingsMap[*vhostAliasEntity.ParentHostname]
		if !mapExists {
			slog.Debug(
				"AliasVirtualHostWithUnknownParent",
				slog.String("method", "SslCertificateWatchdog"),
				slog.String("aliasHostname", vhostAliasEntity.Hostname.String()),
				slog.String("parentHostname", vhostAliasEntity.ParentHostname.String()),
			)
			continue
		}

		vhostSettings.AliasHostnames = append(vhostSettings.AliasHostnames, vhostAliasEntity.Hostname)
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
				slog.String("method", "SslCertificateWatchdog"),
				slog.String("sslPairId", sslPairEntity.Id.String()),
				slog.String("sslPairHostname", sslPairEntity.MainVirtualHostHostname.String()),
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
				slog.String("method", "SslCertificateWatchdog"),
				slog.String("hostname", vhostSettings.VirtualHost.Hostname.String()),
				slog.String("sslPairId", vhostSettings.SslPair.Id.String()),
				slog.String("sslPairHostname", vhostSettings.SslPair.MainVirtualHostHostname.String()),
			)
		}
	}
}
