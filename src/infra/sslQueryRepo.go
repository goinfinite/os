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

func (repo SslQueryRepo) GetHttpdVhostsConfig() (map[string]string, error) {
	httpdVhostsConfig := make(map[string]string)
	httpdVhostsConfigOutput, err := infraHelper.RunCmd(
		"sed", "-n", "/virtualhost/, /}/p", olsHttpdConfigPath,
	)
	if err != nil {
		return httpdVhostsConfig, err
	}

	httpdVhostsConfigSlice := strings.SplitAfter(httpdVhostsConfigOutput, "}\nvirtualhost")

	for _, httpdVhostConfigStr := range httpdVhostsConfigSlice {
		httpdVhostConfigVirtualHostExpression := "virtualhost\\s*(.*)\\s{"
		httpdVhostConfigVirtualHostRegex := regexp.MustCompile(httpdVhostConfigVirtualHostExpression)
		httpdVhostConfigVirtualHostMatch := httpdVhostConfigVirtualHostRegex.FindStringSubmatch(httpdVhostConfigStr)[1]

		httpdVhostConfigFilePathExpression := "configFile\\s*(.*)"
		httpdVhostConfigFilePathRegex := regexp.MustCompile(httpdVhostConfigFilePathExpression)
		httpdVhostConfigFilePathMatch := httpdVhostConfigFilePathRegex.FindStringSubmatch(httpdVhostConfigStr)[1]

		httpdVhostsConfig[httpdVhostConfigVirtualHostMatch] = httpdVhostConfigFilePathMatch
	}

	return httpdVhostsConfig, nil
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
	for _, sslCertificate := range sslCertificates {
		if !sslCertificate.IsCA {
			certificate = sslCertificate
			continue
		}

		chainCertificates = append(chainCertificates, sslCertificate)
	}

	id, err := valueObject.NewSslSerialNumber(certificate.SerialNumber.String())
	if err != nil {
		return ssl, err
	}

	return entity.NewSslPair(
		id,
		hostname,
		certificate,
		privateKey,
		chainCertificates,
	), nil
}

func (repo SslQueryRepo) GetVhostConfigFilePath(vhost valueObject.Fqdn) (string, error) {
	httpdVhostsConfig, err := repo.GetHttpdVhostsConfig()
	if err != nil {
		return "", err
	}

	for virtualHost, configFilePath := range httpdVhostsConfig {
		if vhost.String() != virtualHost {
			continue
		}

		return configFilePath, nil
	}

	return "", errors.New("VhostNotFound")
}

func (repo SslQueryRepo) GetSslPairs() ([]entity.SslPair, error) {
	var sslPairs []entity.SslPair
	httpdVhostsConfig, err := repo.GetHttpdVhostsConfig()
	if err != nil {
		return []entity.SslPair{}, err
	}

	for virtualHost, configFilePath := range httpdVhostsConfig {
		vhostConfigOutputStr, err := infraHelper.RunCmd(
			"sed", "-n", "/vhssl/, /}/p", configFilePath,
		)
		if err != nil {
			return []entity.SslPair{}, err
		}

		if len(vhostConfigOutputStr) < 1 {
			return []entity.SslPair{}, nil
		}

		vhostConfigKeyFileExpression := "keyFile\\s*(.*)"
		vhostConfigKeyFileRegex := regexp.MustCompile(vhostConfigKeyFileExpression)
		vhostConfigKeyFileMatch := vhostConfigKeyFileRegex.FindStringSubmatch(vhostConfigOutputStr)[1]
		privateKeyBytesOutput, err := os.ReadFile(vhostConfigKeyFileMatch)
		if err != nil {
			log.Printf("FailedToOpenHttpdFile: %v", err)
			return []entity.SslPair{}, errors.New("FailedToOpenHttpdFile")
		}
		privateKeyOutputStr := string(privateKeyBytesOutput)

		vhostConfigCertFileExpression := "certFile\\s*(.*)"
		vhostConfigCertFileRegex := regexp.MustCompile(vhostConfigCertFileExpression)
		vhostConfigCertFileMatch := vhostConfigCertFileRegex.FindStringSubmatch(vhostConfigOutputStr)[1]
		certFileBytesOutput, err := os.ReadFile(vhostConfigCertFileMatch)
		if err != nil {
			log.Printf("FailedToOpenVhconfFile: %v", err)
			return []entity.SslPair{}, errors.New("FailedToOpenVhconfFile")
		}
		certFileOutputStr := string(certFileBytesOutput)

		ssl, err := repo.SslFactory(virtualHost, privateKeyOutputStr, certFileOutputStr)
		if err != nil {
			return []entity.SslPair{}, err
		}

		sslPairs = append(sslPairs, ssl)
	}

	return sslPairs, nil
}

func (repo SslQueryRepo) GetSslPairBySerialNumber(sslSerialNumber valueObject.SslSerialNumber) (entity.SslPair, error) {
	sslPairs, err := repo.GetSslPairs()
	if err != nil {
		return entity.SslPair{}, err
	}

	if len(sslPairs) < 1 {
		return entity.SslPair{}, errors.New("SslPairNotFound")
	}

	for _, ssl := range sslPairs {
		if ssl.SerialNumber.String() != sslSerialNumber.String() {
			continue
		}

		return ssl, nil
	}

	return entity.SslPair{}, errors.New("SslPairNotFound")
}
