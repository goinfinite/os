package sslInfra

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log/slog"
	"slices"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type SslQueryRepo struct{}

func NewSslQueryRepo() *SslQueryRepo {
	return &SslQueryRepo{}
}

func (repo *SslQueryRepo) sslCertificatesFactory(
	sslCertsContent valueObject.SslCertificateContent,
) (entity.SslCertificate, []entity.SslCertificate, error) {
	mainCert := entity.SslCertificate{}
	chainedCerts := []entity.SslCertificate{}

	rawSslCertsContent := strings.SplitAfter(
		sslCertsContent.String(), "-----END CERTIFICATE-----\n",
	)
	for _, rawSslCertContent := range rawSslCertsContent {
		if len(rawSslCertContent) == 0 {
			continue
		}

		certificateContent, err := valueObject.NewSslCertificateContent(rawSslCertContent)
		if err != nil {
			return mainCert, chainedCerts, err
		}

		certificate, err := entity.NewSslCertificate(certificateContent)
		if err != nil {
			return mainCert, chainedCerts, err
		}

		if certificate.IsIntermediary {
			chainedCerts = append(chainedCerts, certificate)
			continue
		}

		mainCert = certificate
	}

	if mainCert.CertificateContent.String() == "" {
		return mainCert, chainedCerts, errors.New("MainCertNotFound")
	}

	return mainCert, chainedCerts, nil
}

func (repo *SslQueryRepo) sslPairFactory(
	crtFilePath valueObject.UnixFilePath,
) (sslPairEntity entity.SslPair, err error) {
	crtKeyFilePath := crtFilePath.ReadWithoutExtension().String() + ".key"
	crtKeyContentStr, err := infraHelper.ReadFileContent(crtKeyFilePath)
	if err != nil {
		return sslPairEntity, errors.New("OpenCertKeyFileError: " + err.Error())
	}
	privateKey, err := valueObject.NewSslPrivateKey(crtKeyContentStr)
	if err != nil {
		return sslPairEntity, err
	}

	crtFileContentStr, err := infraHelper.ReadFileContent(crtFilePath.String())
	if err != nil {
		return sslPairEntity, errors.New("OpenCertFileError: " + err.Error())
	}
	certificate, err := valueObject.NewSslCertificateContent(crtFileContentStr)
	if err != nil {
		return sslPairEntity, err
	}

	mainCert, chainedCerts, err := repo.sslCertificatesFactory(certificate)
	if err != nil {
		return sslPairEntity, errors.New("CertsFactoryError: " + err.Error())
	}

	var chainCertificatesContent []valueObject.SslCertificateContent
	for _, sslChainCertificate := range chainedCerts {
		chainCertificatesContent = append(
			chainCertificatesContent, sslChainCertificate.CertificateContent,
		)
	}

	sslPairHashId, err := valueObject.NewSslPairIdFromSslPairContent(
		mainCert.CertificateContent, chainCertificatesContent, privateKey,
	)
	if err != nil {
		return sslPairEntity, err
	}

	crtFileNameWithoutExt := crtFilePath.ReadFileNameWithoutExtension()
	virtualHostHostname, err := valueObject.NewFqdn(crtFileNameWithoutExt.String())
	if err != nil {
		if mainCert.CommonName == nil {
			return sslPairEntity, errors.New("VirtualHostHostnameError: " + err.Error())
		}

		mainCertSslHostname, err := valueObject.NewFqdn(mainCert.CommonName.String())
		if err != nil {
			return sslPairEntity, errors.New("VirtualHostHostnameFallbackError: " + err.Error())
		}

		virtualHostHostname = mainCertSslHostname
	}

	return entity.NewSslPair(
		sslPairHashId, virtualHostHostname, mainCert, privateKey, chainedCerts,
	), nil
}

func (repo *SslQueryRepo) Read(
	requestDto dto.ReadSslPairsRequest,
) (responseDto dto.ReadSslPairsResponse, err error) {
	sslPairEntities := []entity.SslPair{}

	rawCertFilePaths, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "find " + infraEnvs.PkiConfDir +
			" \\( -type f -o -type l \\) -name *.crt",
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return responseDto, errors.New("FindCertFilesError: " + err.Error())
	}

	for rawCertFilePath := range strings.SplitSeq(rawCertFilePaths, "\n") {
		crtFilePath, err := valueObject.NewUnixFilePath(rawCertFilePath)
		if err != nil {
			slog.Debug("InvalidCertFilePath", slog.String("rawCertFilePath", rawCertFilePath))
			continue
		}

		sslPairEntity, err := repo.sslPairFactory(crtFilePath)
		if err != nil {
			slog.Debug("SslPairFactoryError", slog.String("crtFilePath", crtFilePath.String()))
			continue
		}

		if requestDto.SslPairId != nil {
			if sslPairEntity.Id != *requestDto.SslPairId {
				continue
			}
		}

		if requestDto.VirtualHostHostname != nil {
			if sslPairEntity.VirtualHostHostname != *requestDto.VirtualHostHostname {
				continue
			}
		}

		for _, altName := range requestDto.AltNames {
			if !slices.Contains(sslPairEntity.Certificate.AltNames, altName) {
				continue
			}
		}

		if requestDto.IssuedBeforeAt != nil {
			if sslPairEntity.Certificate.IssuedAt >= *requestDto.IssuedBeforeAt {
				continue
			}
		}
		if requestDto.IssuedAfterAt != nil {
			if sslPairEntity.Certificate.IssuedAt <= *requestDto.IssuedAfterAt {
				continue
			}
		}
		if requestDto.ExpiresBeforeAt != nil {
			if sslPairEntity.Certificate.ExpiresAt >= *requestDto.ExpiresBeforeAt {
				continue
			}
		}
		if requestDto.ExpiresAfterAt != nil {
			if sslPairEntity.Certificate.ExpiresAt <= *requestDto.ExpiresAfterAt {
				continue
			}
		}

		sslPairEntities = append(sslPairEntities, sslPairEntity)
	}

	sortDirectionStr := "asc"
	if requestDto.Pagination.SortDirection != nil {
		sortDirectionStr = requestDto.Pagination.SortDirection.String()
	}

	if requestDto.Pagination.SortBy != nil {
		slices.SortStableFunc(sslPairEntities, func(a, b entity.SslPair) int {
			firstElement := a
			secondElement := b
			if sortDirectionStr != "asc" {
				firstElement = b
				secondElement = a
			}

			switch requestDto.Pagination.SortBy.String() {
			case "id", "pairId", "sslPairId":
				return strings.Compare(
					firstElement.Id.String(), secondElement.Id.String(),
				)
			case "virtualHostHostname", "vhostHostname", "hostname":
				return strings.Compare(
					firstElement.VirtualHostHostname.String(), secondElement.VirtualHostHostname.String(),
				)
			case "issuedAt":
				return strings.Compare(
					firstElement.Certificate.IssuedAt.String(), secondElement.Certificate.IssuedAt.String(),
				)
			case "expiresAt":
				return strings.Compare(
					firstElement.Certificate.ExpiresAt.String(), secondElement.Certificate.ExpiresAt.String(),
				)
			default:
				return 0
			}
		})
	}

	itemsTotal := uint64(len(sslPairEntities))
	pagesTotal := uint32(itemsTotal / uint64(requestDto.Pagination.ItemsPerPage))

	paginationDto := requestDto.Pagination
	paginationDto.ItemsTotal = &itemsTotal
	paginationDto.PagesTotal = &pagesTotal

	return dto.ReadSslPairsResponse{
		Pagination: paginationDto,
		SslPairs:   sslPairEntities,
	}, nil
}

func (repo *SslQueryRepo) ReadFirst(
	requestDto dto.ReadSslPairsRequest,
) (sslPairEntity entity.SslPair, err error) {
	requestDto.Pagination = dto.PaginationSingleItem
	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return sslPairEntity, err
	}

	if len(responseDto.SslPairs) == 0 {
		return sslPairEntity, errors.New("SslPairNotFound")
	}

	return responseDto.SslPairs[0], nil
}

func (repo SslQueryRepo) GetOwnershipValidationHash(
	sslCrtContent valueObject.SslCertificateContent,
) (valueObject.Hash, error) {
	sslCrtContentBytes := []byte(sslCrtContent.String())
	sslCrtContentHash := md5.Sum(sslCrtContentBytes)
	sslCrtContentHashStr := hex.EncodeToString(sslCrtContentHash[:])
	return valueObject.NewHash(sslCrtContentHashStr)
}
