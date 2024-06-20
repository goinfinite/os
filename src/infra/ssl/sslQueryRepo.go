package sslInfra

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	infraData "github.com/speedianet/os/src/infra/infraData"
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

		if certificate.IsIntermediary {
			certificates.ChainedCertificates = append(certificates.ChainedCertificates, certificate)

			continue
		}

		certificates.MainCertificate = certificate
	}

	return certificates, nil
}

func (repo SslQueryRepo) sslPairFactory(
	crtFilePath valueObject.UnixFilePath,
) (entity.SslPair, error) {
	var ssl entity.SslPair

	crtKeyFilePath := crtFilePath.GetWithoutExtension().String() + ".key"
	crtKeyContentStr, err := infraHelper.GetFileContent(crtKeyFilePath)
	if err != nil {
		return ssl, errors.New("FailedToOpenCertKeyFile: " + err.Error())
	}
	privateKey, err := valueObject.NewSslPrivateKey(crtKeyContentStr)
	if err != nil {
		return ssl, err
	}

	crtFileContentStr, err := infraHelper.GetFileContent(crtFilePath.String())
	if err != nil {
		return ssl, errors.New("FailedToOpenCertFile: " + err.Error())
	}
	certificate, err := valueObject.NewSslCertificateContent(crtFileContentStr)
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

	crtFileNameWithoutExt := crtFilePath.GetFileNameWithoutExtension()
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

func (repo SslQueryRepo) Read() ([]entity.SslPair, error) {
	sslPairs := []entity.SslPair{}

	crtFilePathsStr, err := infraHelper.RunCmdWithSubShell(
		"find " + infraData.GlobalConfigs.PkiConfDir +
			" \\( -type f -o -type l \\) -name *.crt",
	)
	if err != nil {
		return sslPairs, errors.New("FailedToGetCertFiles: " + err.Error())
	}

	crtFilePaths := strings.Split(crtFilePathsStr, "\n")

	sslPairIdsVhostsNamesMap := map[valueObject.SslId][]valueObject.Fqdn{}
	for _, crtFilePathStr := range crtFilePaths {
		crtFilePath, err := valueObject.NewUnixFilePath(crtFilePathStr)
		if err != nil {
			log.Printf("%s: %s", err.Error(), crtFilePathStr)
			continue
		}

		sslPair, err := repo.sslPairFactory(crtFilePath)
		if err != nil {
			log.Printf("FailedToReadSslPair (%s): %s", crtFilePath, err.Error())
			continue
		}

		pairMainVhostName := sslPair.VirtualHostsHostnames[0]

		_, pairIdAlreadyExists := sslPairIdsVhostsNamesMap[sslPair.Id]
		if pairIdAlreadyExists {
			sslPairIdsVhostsNamesMap[sslPair.Id] = append(
				sslPairIdsVhostsNamesMap[sslPair.Id],
				pairMainVhostName,
			)
			continue
		}

		sslPairIdsVhostsNamesMap[sslPair.Id] = []valueObject.Fqdn{pairMainVhostName}
		sslPairs = append(sslPairs, sslPair)
	}

	for sslPairIndex, sslPair := range sslPairs {
		correctSslPairsVhostsNames := sslPairIdsVhostsNamesMap[sslPair.Id]
		sslPair.VirtualHostsHostnames = correctSslPairsVhostsNames
		sslPairs[sslPairIndex] = sslPair
	}

	return sslPairs, nil
}

func (repo SslQueryRepo) ReadById(sslId valueObject.SslId) (entity.SslPair, error) {
	sslPairs, err := repo.Read()
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

func (repo SslQueryRepo) GetOwnershipValidationHash(
	sslCrtContent valueObject.SslCertificateContent,
) (valueObject.Hash, error) {
	sslCrtContentBytes := []byte(sslCrtContent.String())
	sslCrtContentHash := md5.Sum(sslCrtContentBytes)
	sslCrtContentHashStr := hex.EncodeToString(sslCrtContentHash[:])
	return valueObject.NewHash(sslCrtContentHashStr)
}
