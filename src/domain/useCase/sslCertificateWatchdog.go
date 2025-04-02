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

type VirtualHostSslPair struct {
	VirtualHost entity.VirtualHost
	SslPair     entity.SslPair
}

func (uc *SslCertificateWatchdog) vhostSslPairFactory() map[valueObject.Fqdn]*VirtualHostSslPair {
	vhostHostnameSslPairMap := map[valueObject.Fqdn]*VirtualHostSslPair{}

	vhostReadResponse, err := uc.vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
		Pagination: dto.PaginationUnpaginated,
	})
	if err != nil {
		slog.Error(
			"ReadVirtualHostInfraError",
			slog.String("error", err.Error()),
			slog.String("method", "SslCertificateWatchdog"),
		)
		return vhostHostnameSslPairMap
	}
	for _, vhostEntity := range vhostReadResponse.VirtualHosts {
		vhostHostnameSslPairMap[vhostEntity.Hostname] = &VirtualHostSslPair{
			VirtualHost: vhostEntity,
			SslPair:     entity.SslPair{},
		}
	}
	if len(vhostHostnameSslPairMap) == 0 {
		return vhostHostnameSslPairMap
	}

	sslPairEntities, err := uc.sslQueryRepo.Read()
	if err != nil {
		slog.Error(
			"ReadSslPairInfraError",
			slog.String("error", err.Error()),
			slog.String("method", "SslCertificateWatchdog"),
		)
		return vhostHostnameSslPairMap
	}

	for _, sslPairEntity := range sslPairEntities {
		_, mapExists := vhostHostnameSslPairMap[sslPairEntity.MainVirtualHostHostname]
		if !mapExists {
			slog.Debug(
				"SslPairWithUnknownVirtualHost",
				slog.String("method", "SslCertificateWatchdog"),
				slog.String("sslPairId", sslPairEntity.Id.String()),
				slog.String("sslPairHostname", sslPairEntity.MainVirtualHostHostname.String()),
			)
			continue
		}

		vhostHostnameSslPairMap[sslPairEntity.MainVirtualHostHostname].SslPair = sslPairEntity
	}

	return vhostHostnameSslPairMap
}

func (uc *SslCertificateWatchdog) Execute() {
	vhostHostnameSslPairMap := uc.vhostSslPairFactory()

	skipSubdomainRegex := regexp.MustCompile(`^[^\.]+\.`)
	for _, vhostPair := range vhostHostnameSslPairMap {
		shouldRenewCert := vhostPair.SslPair.IsPubliclyTrusted()

		certAltNamesStrMap := map[string]interface{}{}
		for _, altName := range vhostPair.SslPair.Certificate.AltNames {
			certAltNamesStrMap[altName.String()] = nil
		}

		missingAltNames := []valueObject.Fqdn{}
		for _, aliasHostname := range vhostPair.VirtualHost.AliasesHostnames {
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
			if !shouldRenewCert {
				slog.Debug(
					"SslPairPubliclyTrustedButMissingAltNames",
					slog.String("method", "SslCertificateWatchdog"),
					slog.String("hostname", vhostPair.VirtualHost.Hostname.String()),
					slog.String("sslPairId", vhostPair.SslPair.Id.String()),
					slog.String("sslPairHostname", vhostPair.SslPair.MainVirtualHostHostname.String()),
					slog.Any("currentAltNames", vhostPair.SslPair.Certificate.AltNames),
					slog.Any("missingAltNames", missingAltNames),
				)
			}
			shouldRenewCert = true
		}

		if !shouldRenewCert {
			continue
		}

		err := uc.sslCmdRepo.CreatePubliclyTrusted(dto.NewCreatePubliclyTrustedSslPair(
			vhostPair.VirtualHost.Hostname, vhostPair.VirtualHost.AliasesHostnames,
			valueObject.AccountIdSystem, valueObject.IpAddressSystem,
		))
		if err != nil {
			slog.Debug(
				"CreatePubliclyTrustedSslPairError",
				slog.String("error", err.Error()),
				slog.String("method", "SslCertificateWatchdog"),
				slog.String("hostname", vhostPair.VirtualHost.Hostname.String()),
				slog.String("sslPairId", vhostPair.SslPair.Id.String()),
				slog.String("sslPairHostname", vhostPair.SslPair.MainVirtualHostHostname.String()),
			)
		}
	}
}
