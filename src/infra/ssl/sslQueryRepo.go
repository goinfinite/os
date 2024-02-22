package sslInfra

import (
	"errors"
	"log"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	filesInfra "github.com/speedianet/os/src/infra/files"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type SslQueryRepo struct{}

type SslCertificates struct {
	MainCertificate     entity.SslCertificate
	ChainedCertificates []entity.SslCertificate
}

func (repo SslQueryRepo) sslCertificatesFactory(
	sslCertContent valueObject.SslCertificateContent,
) (SslCertificates, error) {
	var certificates SslCertificates

	sslCertContentSlice := strings.SplitAfter(
		sslCertContent.String(),
		"-----END CERTIFICATE-----\n",
	)
	for _, sslCertContentStr := range sslCertContentSlice {
		if len(sslCertContentStr) == 0 {
			continue
		}

		certificateContent, err := valueObject.NewSslCertificateContent(sslCertContentStr)
		if err != nil {
			return certificates, err
		}

		certificate, err := entity.NewSslCertificate(certificateContent)
		if err != nil {
			return certificates, err
		}

		if !certificate.IsCA {
			certificates.MainCertificate = certificate
			continue
		}

		certificates.ChainedCertificates = append(certificates.ChainedCertificates, certificate)
	}

	return certificates, nil
}

func (repo SslQueryRepo) sslPairFactory(
	crtFile entity.UnixFile,
) (entity.SslPair, error) {
	var ssl entity.SslPair

	vhostCrtKeyFilePath := crtFile.Path.GetWithoutExtension().String() + ".key"
	vhostCrtKeyContentStr, err := infraHelper.GetFileContent(vhostCrtKeyFilePath)
	if err != nil {
		return ssl, errors.New("FailedToOpenCertKeyFile: " + err.Error())
	}
	privateKey, err := valueObject.NewSslPrivateKey(vhostCrtKeyContentStr)
	if err != nil {
		return ssl, err
	}

	vhostCrtFilePath := crtFile.Path.String()
	vhostCrtFileContentStr, err := infraHelper.GetFileContent(vhostCrtFilePath)
	if err != nil {
		return ssl, errors.New("FailedToOpenCertFile: " + err.Error())
	}
	certificate, err := valueObject.NewSslCertificateContent(vhostCrtFileContentStr)
	if err != nil {
		return ssl, err
	}

	sslCertificates, err := repo.sslCertificatesFactory(certificate)
	if err != nil {
		return ssl, errors.New("FailedToGetMainAndChainedCerts: " + err.Error())
	}

	mainCertificate := sslCertificates.MainCertificate
	chainCertificates := sslCertificates.ChainedCertificates

	var chainCertificatesContent []valueObject.SslCertificateContent
	for _, sslChainCertificate := range chainCertificates {
		chainCertificatesContent = append(
			chainCertificatesContent,
			sslChainCertificate.CertificateContent,
		)
	}

	hashId, err := valueObject.NewSslIdFromSslPairContent(
		mainCertificate.CertificateContent,
		chainCertificatesContent,
		privateKey,
	)
	if err != nil {
		return ssl, err
	}

	crtFileNameWithoutExt := crtFile.Path.GetFileNameWithoutExtension()
	vhost, err := valueObject.NewFqdn(crtFileNameWithoutExt.String())
	if err != nil {
		return ssl, err
	}

	return entity.NewSslPair(
		hashId,
		[]valueObject.Fqdn{vhost},
		mainCertificate,
		privateKey,
		chainCertificates,
	), nil
}

func (repo SslQueryRepo) GetSslPairs() ([]entity.SslPair, error) {
	sslPairs := []entity.SslPair{}

	pkiDirPath, _ := valueObject.NewUnixFilePath("/app/conf/pki")

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	files, err := filesQueryRepo.Get(pkiDirPath)
	if err != nil {
		return sslPairs, errors.New("FailedToGetFiles: " + err.Error())
	}

	crtFiles := []entity.UnixFile{}
	for _, file := range files {
		if file.MimeType.IsDir() {
			continue
		}

		if file.Extension.String() == "crt" {
			crtFiles = append(crtFiles, file)
		}
	}

	repeatedSslPairIdsWithVhosts := map[valueObject.SslId][]valueObject.Fqdn{}
	for _, crtFile := range crtFiles {
		sslPair, err := repo.sslPairFactory(crtFile)
		if err != nil {
			log.Printf("FailedToGetSslPair (%s): %s", crtFile.Path, err.Error())
		}

		_, sslPairIdAlreadySavedInMap := repeatedSslPairIdsWithVhosts[sslPair.Id]
		if !sslPairIdAlreadySavedInMap {
			repeatedSslPairIdsWithVhosts[sslPair.Id] = []valueObject.Fqdn{}

			sslPairs = append(sslPairs, sslPair)

			continue
		}

		uniqueVhost := sslPair.VirtualHosts[0]
		repeatedSslPairIdsWithVhosts[sslPair.Id] = append(
			repeatedSslPairIdsWithVhosts[sslPair.Id],
			uniqueVhost,
		)
	}

	for sslPairIndex, sslPair := range sslPairs {
		_, sslPairIdExistsInMap := repeatedSslPairIdsWithVhosts[sslPair.Id]
		if !sslPairIdExistsInMap {
			continue
		}

		sslPair.VirtualHosts = append(
			sslPair.VirtualHosts,
			repeatedSslPairIdsWithVhosts[sslPair.Id]...,
		)

		sslPairs[sslPairIndex] = sslPair
	}

	return sslPairs, nil
}

func (repo SslQueryRepo) GetSslPairById(sslId valueObject.SslId) (entity.SslPair, error) {
	sslPairs, err := repo.GetSslPairs()
	if err != nil {
		return entity.SslPair{}, err
	}

	if len(sslPairs) < 1 {
		return entity.SslPair{}, errors.New("SslPairNotFound")
	}

	for _, ssl := range sslPairs {
		if ssl.Id.String() != sslId.String() {
			continue
		}

		return ssl, nil
	}

	return entity.SslPair{}, errors.New("SslPairNotFound")
}
