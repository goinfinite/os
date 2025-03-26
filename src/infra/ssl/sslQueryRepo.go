package sslInfra

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"slices"
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
	crtKeyContentStr, err := infraHelper.GetFileContent(crtKeyFilePath)
	if err != nil {
		return sslPairEntity, errors.New("OpenCertKeyFileError: " + err.Error())
	}
	privateKey, err := valueObject.NewSslPrivateKey(crtKeyContentStr)
	if err != nil {
		return sslPairEntity, err
	}

	crtFileContentStr, err := infraHelper.GetFileContent(crtFilePath.String())
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
		sslPairHashId, mainVirtualHostHostname, []valueObject.Fqdn{mainVirtualHostHostname},
		mainCert, privateKey, chainedCerts,
	), nil
}

func (repo SslQueryRepo) Read() ([]entity.SslPair, error) {
	sslPairs := []entity.SslPair{}

	crtFilePathsStr, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "find " + infraEnvs.PkiConfDir +
			" \\( -type f -o -type l \\) -name *.crt",
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return sslPairs, errors.New("FailedToGetCertFiles: " + err.Error())
	}

	crtFilePaths := strings.Split(crtFilePathsStr, "\n")

	sslPairIdsVhostsNamesMap := map[valueObject.SslPairId][]valueObject.Fqdn{}
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

func (repo SslQueryRepo) ReadByVhostHostname(
	sslPairVhostHostname valueObject.Fqdn,
) (entity.SslPair, error) {
	sslPairs, err := repo.Read()
	if err != nil {
		return entity.SslPair{}, err
	}

	if len(sslPairs) < 1 {
		return entity.SslPair{}, errors.New("SslPairNotFound")
	}

	for _, ssl := range sslPairs {
		if !slices.Contains(ssl.VirtualHostsHostnames, sslPairVhostHostname) {
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
