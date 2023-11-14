package infra

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

const olsHttpdConfigPath = "/usr/local/lsws/conf/httpd_config.conf"

type SslQueryRepo struct{}

type SslCertificates struct {
	MainCertificate     entity.SslCertificate
	ChainedCertificates []entity.SslCertificate
}

func (repo SslQueryRepo) GetVhosts() ([]valueObject.Fqdn, error) {
	httpdContent, err := infraHelper.GetFileContent(olsHttpdConfigPath)
	if err != nil {
		return []valueObject.Fqdn{}, err
	}

	vhostsExpression := `virtualhost\s*(.*) {`
	vhostsRegex := regexp.MustCompile(vhostsExpression)
	vhostsMatch := vhostsRegex.FindAllStringSubmatch(httpdContent, -1)
	if len(vhostsMatch) < 1 {
		return []valueObject.Fqdn{}, err
	}

	httpdVhosts := []valueObject.Fqdn{}
	for _, vhostMatchStr := range vhostsMatch {
		if len(vhostMatchStr) < 2 {
			continue
		}

		vhostStr := vhostMatchStr[1]
		vhost, err := valueObject.NewFqdn(vhostStr)
		if err != nil {
			log.Printf("UnableToGetVhost (%v): %v", vhostStr, err)
			continue
		}
		httpdVhosts = append(httpdVhosts, vhost)
	}

	return httpdVhosts, nil
}

func (repo SslQueryRepo) GetVhostConfigFilePath(
	vhost valueObject.Fqdn,
) (valueObject.UnixFilePath, error) {
	var vhostConfigFilePath valueObject.UnixFilePath
	httpdContent, err := infraHelper.GetFileContent(olsHttpdConfigPath)
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

	_, err := repo.GetVhostConfigFilePath(sslHostname)
	if err != nil {
		return ssl, err
	}

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
	httpdVhosts, err := repo.GetVhosts()
	if err != nil {
		return []entity.SslPair{}, err
	}

	for _, vhost := range httpdVhosts {
		vhostConfigFilePath, err := repo.GetVhostConfigFilePath(vhost)
		if err != nil {
			return []entity.SslPair{}, err
		}

		vhostConfigContentStr, err := infraHelper.GetFileContent(vhostConfigFilePath.String())
		if err != nil {
			return []entity.SslPair{}, err
		}

		if len(vhostConfigContentStr) < 1 {
			return []entity.SslPair{}, nil
		}

		vhostConfigKeyFileExpression := `keyFile\s*(.*)`
		vhostConfigKeyFileMatch, err := infraHelper.GetRegexFirstGroup(vhostConfigContentStr, vhostConfigKeyFileExpression)
		if err != nil {
			return []entity.SslPair{}, nil
		}
		privateKeyContentStr, err := infraHelper.GetFileContent(vhostConfigKeyFileMatch)
		if err != nil {
			log.Printf("FailedToOpenHttpdFile: %v", err)
			return []entity.SslPair{}, errors.New("FailedToOpenHttpdFile")
		}
		privateKey, err := valueObject.NewSslPrivateKey(privateKeyContentStr)
		if err != nil {
			return []entity.SslPair{}, nil
		}

		vhostConfigCertFileExpression := `certFile\s*(.*)`
		vhostConfigCertFileMatch, err := infraHelper.GetRegexFirstGroup(vhostConfigContentStr, vhostConfigCertFileExpression)
		if err != nil {
			return []entity.SslPair{}, nil
		}
		certFileContentStr, err := infraHelper.GetFileContent(vhostConfigCertFileMatch)
		if err != nil {
			log.Printf("FailedToOpenVhconfFile: %v", err)
			return []entity.SslPair{}, errors.New("FailedToOpenVhconfFile")
		}
		certificate, err := valueObject.NewSslCertificateContent(certFileContentStr)
		if err != nil {
			return []entity.SslPair{}, nil
		}

		sslCertificates, err := repo.SslCertificatesFactory(certificate)
		if err != nil {
			return []entity.SslPair{}, err
		}

		ssl, err := repo.SslPairFactory(vhost, privateKey, sslCertificates)
		if err != nil {
			return []entity.SslPair{}, err
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
