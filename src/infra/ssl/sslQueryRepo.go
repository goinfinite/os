package sslInfra

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log/slog"
	"strings"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type SslQueryRepo struct{}

func (repo SslQueryRepo) sslCertificatesFactory(
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

func (repo SslQueryRepo) sslPairFactory(
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
	mainVirtualHostHostname, err := valueObject.NewFqdn(crtFileNameWithoutExt.String())
	if err != nil {
		if mainCert.CommonName == nil {
			return sslPairEntity, errors.New("MainVirtualHostHostnameError: " + err.Error())
		}

		mainCertSslHostname, err := valueObject.NewFqdn(mainCert.CommonName.String())
		if err != nil {
			return sslPairEntity, errors.New("MainVirtualHostHostnameFallbackError: " + err.Error())
		}

		mainVirtualHostHostname = mainCertSslHostname
	}

	return entity.NewSslPair(
		sslPairHashId, mainVirtualHostHostname, mainCert, privateKey, chainedCerts,
	), nil
}

func (repo SslQueryRepo) Read() ([]entity.SslPair, error) {
	sslPairEntities := []entity.SslPair{}

	rawCertFilePaths, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "find " + infraEnvs.PkiConfDir +
			" \\( -type f -o -type l \\) -name *.crt",
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return sslPairEntities, errors.New("FindCertFilesError: " + err.Error())
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

		sslPairEntities = append(sslPairEntities, sslPairEntity)
	}

	return sslPairEntities, nil
}

func (repo SslQueryRepo) ReadById(
	sslPairId valueObject.SslPairId,
) (entity.SslPair, error) {
	sslPairs, err := repo.Read()
	if err != nil {
		return entity.SslPair{}, err
	}

	if len(sslPairs) < 1 {
		return entity.SslPair{}, errors.New("SslPairNotFound")
	}

	for _, ssl := range sslPairs {
		if ssl.Id.String() != sslPairId.String() {
			continue
		}

		return ssl, nil
	}

	return entity.SslPair{}, errors.New("SslPairNotFound")
}

func (repo SslQueryRepo) GetOwnershipValidationHash(
	sslCrtContent valueObject.SslCertificateContent,
) (valueObject.Hash, error) {
	sslCrtContentBytes := []byte(sslCrtContent.String())
	sslCrtContentHash := md5.Sum(sslCrtContentBytes)
	sslCrtContentHashStr := hex.EncodeToString(sslCrtContentHash[:])
	return valueObject.NewHash(sslCrtContentHashStr)
}
