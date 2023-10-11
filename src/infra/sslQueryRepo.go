package infra

import (
	"errors"
	"strings"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

const olsHttpdConfigPath = "/usr/local/lsws/conf/httpd_config.conf"

type SslQueryRepo struct{}

type HttpdVhostConfig struct {
	VirtualHost string
	FilePath    string
}

type VhostConfig struct {
	FilePath    string
	FileContent string
}

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

func (repo SslQueryRepo) SslFactory(
	sslHostname string,
	sslPrivateKey string,
	sslCertContent string,
) (entity.SslPair, error) {
	var ssl entity.SslPair

	hostname, err := valueObject.NewFqdn(sslHostname)
	if err != nil {
		return ssl, errors.New("SslHostnameError")
	}
	_, err = repo.GetVhostConfig(hostname.String())
	if err != nil {
		return ssl, err
	}

	privateKey, err := entity.NewSslPrivateKey(sslPrivateKey)
	if err != nil {
		return ssl, err
	}

	sslCertificates, err := repo.splitSslCertificate(sslCertContent)
	if err != nil {
		return ssl, err
	}

	var certificate entity.SslCertificate
	var chainCertificates []entity.SslCertificate
	for _, sslCertificate := range sslCertificates {
		if sslCertificate.IsCA {
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

func (repo SslQueryRepo) GetHttpdVhostsConfig() ([]HttpdVhostConfig, error) {
	var httpdVhostsConfig []HttpdVhostConfig
	httpdVhostsConfigOutput, err := infraHelper.RunCmd(
		"sed", "-n", "/virtualhost/, /}/p", olsHttpdConfigPath,
	)
	if err != nil {
		return []HttpdVhostConfig{}, err
	}

	httpdVhostsConfigSlice := strings.SplitAfter(httpdVhostsConfigOutput, "}\nvirtualhost")

	for _, httpdVhostConfigStr := range httpdVhostsConfigSlice {
		httpdVhostConfigGroups := infraHelper.GetRegexNamedGroups(httpdVhostConfigStr, "(?:virtualhost )(?P<virtualHost>.*) {\n\\s*vhRoot\\s*.*\n\\s*(?:configFile\\s*)(?P<configFile>.*)")
		httpdVhostConfigVirtualHost := httpdVhostConfigGroups["virtualHost"]
		httpdVhostConfigFilePath := httpdVhostConfigGroups["configFile"]

		httpdVhostsConfig = append(httpdVhostsConfig, HttpdVhostConfig{
			VirtualHost: httpdVhostConfigVirtualHost,
			FilePath:    httpdVhostConfigFilePath,
		})
	}

	return httpdVhostsConfig, nil
}

func (repo SslQueryRepo) GetVhostConfig(vhost string) (VhostConfig, error) {
	vhostConfigFilePath := ""
	vhostConfigFileContent := ""

	httpdVhostsConfigs, err := repo.GetHttpdVhostsConfig()
	if err != nil {
		return VhostConfig{}, err
	}

	for _, httpdVhostConfig := range httpdVhostsConfigs {
		if vhost != httpdVhostConfig.VirtualHost {
			continue
		}

		vhostConfigFilePath = httpdVhostConfig.FilePath
		vhostConfigFileContent, err = infraHelper.RunCmd("cat", vhostConfigFilePath)
		if err != nil {
			return VhostConfig{}, err
		}
		break
	}

	if len(vhostConfigFilePath) == 0 {
		return VhostConfig{}, errors.New("VhostNotFound")
	}

	return VhostConfig{
		FilePath:    vhostConfigFilePath,
		FileContent: vhostConfigFileContent,
	}, nil
}

func (repo SslQueryRepo) GetSslPairs() ([]entity.SslPair, error) {
	var sslPairs []entity.SslPair
	httpdVhostsConfig, err := repo.GetHttpdVhostsConfig()
	if err != nil {
		return []entity.SslPair{}, err
	}

	for _, httpdVhostConfig := range httpdVhostsConfig {
		vhostConfigOutput, err := infraHelper.RunCmd(
			"sed", "-n", "/vhssl/, /}/p", httpdVhostConfig.FilePath,
		)
		if err != nil {
			return []entity.SslPair{}, err
		}

		if len(vhostConfigOutput) < 1 {
			return []entity.SslPair{}, nil
		}

		vhostConfigGroups := infraHelper.GetRegexNamedGroups(vhostConfigOutput, "(?:keyFile\\s*(?P<keyFile>.*))?\n\\s*(?:certFile\\s*(?P<certFile>.*))\n\\s*(?:certChain\\s*(?P<certChain>.*))\n\\s*(?:(?:CACertPath\\s*(?P<CACertPath>.*))\n\\s*)?(?:(?:CACertFile\\s*(?P<CACertFile>.*))?)")
		privateKeyOutput, err := infraHelper.RunCmd("cat", vhostConfigGroups["keyFile"])
		if err != nil {
			return []entity.SslPair{}, err
		}

		certFileOutput, err := infraHelper.RunCmd("cat", vhostConfigGroups["certFile"])
		if err != nil {
			return []entity.SslPair{}, err
		}

		ssl, err := repo.SslFactory(httpdVhostConfig.VirtualHost, privateKeyOutput, certFileOutput)
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
