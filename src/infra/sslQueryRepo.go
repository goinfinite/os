package infra

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

const olsHttpdConfigPath = "/usr/local/lsws/conf/httpd_config.conf"

type SslQueryRepo struct{}

func (repo SslQueryRepo) splitSslCertificate(
	sslCertContent string,
) ([]entity.SslCertificate, error) {
	certificates := []entity.SslCertificate{}

	sslCertContentSlice := strings.SplitAfter(sslCertContent, "-----END CERTIFICATE-----\n")
	for _, sslCertContentStr := range sslCertContentSlice {
		certificate, err := entity.NewSslCertificate(sslCertContentStr)
		if err != nil {
			return certificates, err
		}
		certificates = append(certificates, certificate)
	}

	return certificates, nil
}

func (repo SslQueryRepo) GetVhosts() ([]valueObject.Fqdn, error) {
	httpdContent, err := infraHelper.ReadFile(olsHttpdConfigPath)
	if err != nil {
		return []valueObject.Fqdn{}, err
	}

	vhostsExpression := "virtualhost\\s*(.*) {"
	vhostsRegex := regexp.MustCompile(vhostsExpression)
	vhostsMatch := vhostsRegex.FindAllStringSubmatch(httpdContent, -1)
	if len(vhostsMatch) < 1 {
		return []valueObject.Fqdn{}, err
	}

	httpdVhosts := []valueObject.Fqdn{}
	for _, vhostMatchStr := range vhostsMatch {
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

func (repo SslQueryRepo) SslFactory(
	sslHostname string,
	sslPrivateKeyStr string,
	sslCertStr string,
) (entity.SslPair, error) {
	var ssl entity.SslPair

	hostname, err := valueObject.NewFqdn(sslHostname)
	if err != nil {
		return ssl, errors.New("SslHostnameError")
	}
	_, err = repo.GetVhostConfigFilePath(hostname)
	if err != nil {
		return ssl, err
	}

	privateKey, err := valueObject.NewSslPrivateKey(sslPrivateKeyStr)
	if err != nil {
		return ssl, err
	}

	sslCertificates, err := repo.splitSslCertificate(sslCertStr)
	if err != nil {
		return ssl, err
	}

	var certificate entity.SslCertificate
	var chainCertificates []entity.SslCertificate
	var chainCertificatesContent []valueObject.SslCertificateStr
	for _, sslCertificate := range sslCertificates {
		if !sslCertificate.IsCA {
			certificate = sslCertificate
			continue
		}

		chainCertificates = append(chainCertificates, sslCertificate)
		chainCertificatesContent = append(chainCertificatesContent, sslCertificate.Certificate)
	}

	hashId, err := valueObject.NewSslHashIdFromSslPair(
		certificate.Certificate,
		chainCertificatesContent,
		privateKey,
	)
	if err != nil {
		return ssl, err
	}

	return entity.NewSslPair(
		hashId,
		hostname,
		certificate,
		privateKey,
		chainCertificates,
	), nil
}

func (repo SslQueryRepo) GetVhostConfigFilePath(
	vhost valueObject.Fqdn,
) (valueObject.UnixFilePath, error) {
	var vhostConfigFilePath valueObject.UnixFilePath
	httpdContent, err := infraHelper.ReadFile(olsHttpdConfigPath)
	if err != nil {
		return "", err
	}

	vhostConfigFileExpression := "\\s*configFile\\s*(.*)"
	vhostConfigFileMatch, err := infraHelper.GetRegexUniqueMatch(httpdContent, vhostConfigFileExpression)
	if err != nil {
		return "", err
	}

	vhostConfigFilePath, err = valueObject.NewUnixFilePath(vhostConfigFileMatch)
	if err != nil {
		return "", err
	}

	return vhostConfigFilePath, nil
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

		vhostConfigContentStr, err := infraHelper.ReadFile(vhostConfigFilePath.String())
		if err != nil {
			return []entity.SslPair{}, err
		}

		/* TODO: Remover quando implementar o middleware de validação de serviço. */
		if len(vhostConfigContentStr) < 1 {
			return []entity.SslPair{}, nil
		}

		vhostConfigKeyFileExpression := "keyFile\\s*(.*)"
		vhostConfigKeyFileMatch, err := infraHelper.GetRegexUniqueMatch(vhostConfigContentStr, vhostConfigKeyFileExpression)
		if err != nil {
			return []entity.SslPair{}, nil
		}
		privateKeyContentBytes, err := os.ReadFile(vhostConfigKeyFileMatch)
		if err != nil {
			log.Printf("FailedToOpenHttpdFile: %v", err)
			return []entity.SslPair{}, errors.New("FailedToOpenHttpdFile")
		}
		privateKeyContentStr := string(privateKeyContentBytes)

		vhostConfigCertFileExpression := "certFile\\s*(.*)"
		vhostConfigCertFileMatch, err := infraHelper.GetRegexUniqueMatch(vhostConfigContentStr, vhostConfigCertFileExpression)
		if err != nil {
			return []entity.SslPair{}, nil
		}
		certFileContentBytes, err := os.ReadFile(vhostConfigCertFileMatch)
		if err != nil {
			log.Printf("FailedToOpenVhconfFile: %v", err)
			return []entity.SslPair{}, errors.New("FailedToOpenVhconfFile")
		}
		certFileContentStr := string(certFileContentBytes)

		ssl, err := repo.SslFactory(vhost.String(), privateKeyContentStr, certFileContentStr)
		if err != nil {
			return []entity.SslPair{}, err
		}

		sslPairs = append(sslPairs, ssl)
	}

	return sslPairs, nil
}

func (repo SslQueryRepo) GetSslPairByHashId(sslHashId valueObject.SslHashId) (entity.SslPair, error) {
	sslPairs, err := repo.GetSslPairs()
	if err != nil {
		return entity.SslPair{}, err
	}

	if len(sslPairs) < 1 {
		return entity.SslPair{}, errors.New("SslPairNotFound")
	}

	for _, ssl := range sslPairs {
		if ssl.HashId.String() != sslHashId.String() {
			continue
		}

		return ssl, nil
	}

	return entity.SslPair{}, errors.New("SslPairNotFound")
}
