package sslInfra

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
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
	sslVhosts []valueObject.Fqdn,
) (entity.SslPair, error) {
	var ssl entity.SslPair

	firstVhost := sslVhosts[0]
	firstVhostStr := firstVhost.String()

	vhostCertKeyFilePath := pkiConfDir + "/" + firstVhostStr + ".key"
	vhostCertKeyContentStr, err := infraHelper.GetFileContent(vhostCertKeyFilePath)
	if err != nil {
		return ssl, errors.New("FailedToOpenCertKeyFile: " + err.Error())
	}
	privateKey, err := valueObject.NewSslPrivateKey(vhostCertKeyContentStr)
	if err != nil {
		return ssl, err
	}

	vhostCertFilePath := pkiConfDir + "/" + firstVhostStr + ".crt"
	vhostCertFileContentStr, err := infraHelper.GetFileContent(vhostCertFilePath)
	if err != nil {
		return ssl, errors.New("FailedToOpenCertFile: " + err.Error())
	}
	certificate, err := valueObject.NewSslCertificateContent(vhostCertFileContentStr)
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

	return entity.NewSslPair(
		hashId,
		sslVhosts,
		mainCertificate,
		privateKey,
		chainCertificates,
	), nil
}

func (repo SslQueryRepo) GetSslPairs() ([]entity.SslPair, error) {
	sslPairs := []entity.SslPair{}

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhosts, err := vhostQueryRepo.Get()
	if err != nil {
		return sslPairs, errors.New("FailedToGetVhosts")
	}

	certFilePathWithVhosts := map[string][]valueObject.Fqdn{}
	for _, vhost := range vhosts {
		certFilePath := pkiConfDir + "/" + vhost.Hostname.String() + ".crt"

		isSymlink := infraHelper.IsSymlink(certFilePath)
		if isSymlink {
			targetCertFilePath, err := os.Readlink(certFilePath)
			if err != nil {
				log.Printf("FailedToGetTargetCertFilePathFromSymlink: %s", err.Error())
				continue
			}

			certFilePath = targetCertFilePath
		}

		_, certFilePathAlreadyExistsInMap := certFilePathWithVhosts[certFilePath]
		if !certFilePathAlreadyExistsInMap {
			certFilePathWithVhosts[certFilePath] = []valueObject.Fqdn{}
		}

		certFilePathWithVhosts[certFilePath] = append(
			certFilePathWithVhosts[certFilePath],
			vhost.Hostname,
		)
	}

	for certFilePath, vhosts := range certFilePathWithVhosts {
		sslPair, err := repo.sslPairFactory(vhosts)
		if err != nil {
			log.Printf("FailedToGetSslPair (%s): %s", certFilePath, err.Error())
			continue
		}

		sslPairs = append(sslPairs, sslPair)
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
