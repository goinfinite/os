package sslInfra

import (
	"errors"
	"log"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

const configurationsDir = "/app/conf/nginx"

type SslQueryRepo struct{}

type SslCertificates struct {
	MainCertificate     entity.SslCertificate
	ChainedCertificates []entity.SslCertificate
}

// TODO: add "getCertFileContentByRegExp(regExp string, vhostConfFilePath string) (string, error)"

func (repo SslQueryRepo) GetVhostConfigFilePath(
	vhost valueObject.Fqdn,
) (valueObject.UnixFilePath, error) {
	var vhostConfigFilePath valueObject.UnixFilePath
	httpdContent, err := infraHelper.GetFileContent(configurationsDir)
	if err != nil {
		return "", err
	}

	vhostConfigFileExpression := `\s*configFile\s*(.*)`
	vhostConfigFileMatch, err := infraHelper.GetRegexFirstGroup(httpdContent, vhostConfigFileExpression)
	if err != nil {
		return "", err
	}

	vhostConfigFilePath, err = valueObject.NewUnixFilePath(vhostConfigFileMatch)
	if err != nil {
		return "", err
	}

	return vhostConfigFilePath, nil
}

func (repo SslQueryRepo) SslCertificatesFactory(
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

func (repo SslQueryRepo) SslPairFactory(
	sslHostname valueObject.Fqdn,
	sslPrivateKey valueObject.SslPrivateKey,
	sslCertificates SslCertificates,
) (entity.SslPair, error) {
	var ssl entity.SslPair

	certificate := sslCertificates.MainCertificate
	chainCertificates := sslCertificates.ChainedCertificates

	var chainCertificatesContent []valueObject.SslCertificateContent
	for _, sslChainCertificate := range chainCertificates {
		chainCertificatesContent = append(chainCertificatesContent, sslChainCertificate.Certificate)
	}

	hashId, err := valueObject.NewSslIdFromSslPairContent(
		certificate.Certificate,
		chainCertificatesContent,
		sslPrivateKey,
	)
	if err != nil {
		return ssl, err
	}

	return entity.NewSslPair(
		hashId,
		sslHostname,
		certificate,
		sslPrivateKey,
		chainCertificates,
	), nil
}

func (repo SslQueryRepo) GetSslPairs() ([]entity.SslPair, error) {
	var sslPairs []entity.SslPair
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}

	virtualHosts, err := vhostQueryRepo.Get()
	if err != nil {
		return []entity.SslPair{}, err
	}

	for _, vhost := range virtualHosts {
		hostnameStr := vhost.Hostname.String()

		vhostConfigFilePath := configurationsDir + "/" + hostnameStr + ".conf"
		if vhostQueryRepo.IsVirtualHostPrimaryDomain(vhost.Hostname) {
			vhostConfigFilePath = configurationsDir + "/primary.conf"
		}

		vhostConfigContentStr, err := infraHelper.GetFileContent(vhostConfigFilePath)
		if err != nil {
			log.Printf("FailedToOpenVhostConfFile (%s): %s", hostnameStr, err.Error())
			continue
		}

		fileIsEmpty := len(vhostConfigContentStr) < 1
		if fileIsEmpty {
			log.Printf("VirtualHostConfFileIsEmpty (%s)", hostnameStr)
			continue
		}

		vhostCertKeyFileExp := `ssl_certificate_key\s*(.*);`
		vhostCertKeyFilePath, err := infraHelper.GetRegexFirstGroup(
			vhostConfigContentStr,
			vhostCertKeyFileExp,
		)
		if err != nil {
			log.Printf("FailedToGetCertKeyFilePath (%s): %s", hostnameStr, err.Error())
			continue
		}
		vhostCertKeyContentStr, err := infraHelper.GetFileContent(vhostCertKeyFilePath)
		if err != nil {
			log.Printf("FailedToOpenCertKeyFile (%s): %s", hostnameStr, err.Error())
			continue
		}
		privateKey, err := valueObject.NewSslPrivateKey(vhostCertKeyContentStr)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), hostnameStr)
			continue
		}

		vhostCertFileExp := `ssl_certificate\s*(.*);`
		vhostCertFilePath, err := infraHelper.GetRegexFirstGroup(
			vhostConfigContentStr,
			vhostCertFileExp,
		)
		if err != nil {
			log.Printf("FailedToGetCertFilePath (%s): %s", hostnameStr, err.Error())
			continue
		}
		vhostCertFileContentStr, err := infraHelper.GetFileContent(vhostCertFilePath)
		if err != nil {
			log.Printf("FailedToOpenCertFile (%s): %s", hostnameStr, err.Error())
			continue
		}
		certificate, err := valueObject.NewSslCertificateContent(vhostCertFileContentStr)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), hostnameStr)
			continue
		}

		sslCertificates, err := repo.SslCertificatesFactory(certificate)
		if err != nil {
			log.Printf("FailedToGetMainAndChainedCerts (%s): %s", hostnameStr, err.Error())
			continue
		}

		ssl, err := repo.SslPairFactory(vhost.Hostname, privateKey, sslCertificates)
		if err != nil {
			log.Printf("FailedToGetSslPair (%s): %s", hostnameStr, err.Error())
			continue
		}

		sslPairs = append(sslPairs, ssl)
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
