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

func (uc *SslCertificateWatchdog) shouldRenewCert(
	vhostEntity entity.VirtualHost,
	sslPairEntity entity.SslPair,
) bool {
	if len(vhostEntity.AliasesHostnames) == 0 {
		return sslPairEntity.IsPubliclyTrusted()
	}

	certAltNamesStrMap := map[string]interface{}{}
	for _, altName := range sslPairEntity.Certificate.AltNames {
		certAltNamesStrMap[altName.String()] = nil
	}

	skipSubdomainRegex := regexp.MustCompile(`^[^\.]+\.`)
	missingAltNames := []valueObject.Fqdn{}
	for _, aliasHostname := range vhostEntity.AliasesHostnames {
		aliasHostnameStr := aliasHostname.String()
		if _, altNameExists := certAltNamesStrMap[aliasHostnameStr]; altNameExists {
			continue
		}

		likelyParentHostname := skipSubdomainRegex.ReplaceAllString(aliasHostnameStr, "")
		wildcardParentHostname := "*." + likelyParentHostname
		if _, altNameExists := certAltNamesStrMap[wildcardParentHostname]; altNameExists {
			continue
		}

		missingAltNames = append(missingAltNames, aliasHostname)
	}

	if len(missingAltNames) == 0 {
		return sslPairEntity.IsPubliclyTrusted()
	}

	if sslPairEntity.IsPubliclyTrusted() {
		slog.Debug(
			"SslPairPubliclyTrustedButMissingAltNames",
			slog.String("method", "SslCertificateWatchdog"),
			slog.String("hostname", vhostEntity.Hostname.String()),
			slog.Any("currentAltNames", sslPairEntity.Certificate.AltNames),
			slog.Any("missingAltNames", missingAltNames),
		)
	}
	return true
}

func (uc *SslCertificateWatchdog) Execute() {
	vhostReadResponse, err := uc.vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
		Pagination: dto.PaginationUnpaginated,
	})
	if err != nil {
		slog.Error(
			"ReadVirtualHostInfraError",
			slog.String("error", err.Error()),
			slog.String("method", "SslCertificateWatchdog"),
		)
		return
	}

	for _, vhostEntity := range vhostReadResponse.VirtualHosts {
		sslPairEntity, err := uc.sslQueryRepo.ReadFirst(dto.ReadSslPairsRequest{
			VirtualHostHostname: &vhostEntity.Hostname,
		})
		if err != nil {
			slog.Debug(
				"ReadSslPairError",
				slog.String("error", err.Error()),
				slog.String("method", "SslCertificateWatchdog"),
				slog.String("hostname", vhostEntity.Hostname.String()),
			)
			continue
		}

		if !uc.shouldRenewCert(vhostEntity, sslPairEntity) {
			continue
		}

		_, err = uc.sslCmdRepo.CreatePubliclyTrusted(dto.NewCreatePubliclyTrustedSslPair(
			vhostEntity.Hostname, valueObject.AccountIdSystem, valueObject.IpAddressSystem,
		))
		if err != nil {
			slog.Debug(
				"CreatePubliclyTrustedSslPairError",
				slog.String("error", err.Error()),
				slog.String("method", "SslCertificateWatchdog"),
				slog.String("hostname", vhostEntity.Hostname.String()),
			)
		}
	}
}
