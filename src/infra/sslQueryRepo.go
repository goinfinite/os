package infra

import (
	"errors"
	"regexp"
	"strings"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

type SslQueryRepo struct {
	olsHttpdConfigPath string
}

type HttpdVhostConfig struct {
	VirtualHost string
	FilePath    string
}

type VhostConfig struct {
	FilePath    string
	FileContent string
}

func NewSslQueryRepo() *SslQueryRepo {
	return &SslQueryRepo{
		olsHttpdConfigPath: "/usr/local/lsws/conf/httpd_config.conf",
	}
}

func (repo SslQueryRepo) splitSslCertificate(
	sslCertContent string,
) ([]entity.SslPair, error) {
	certificates := []entity.SslPair{}

	sslCertContentSlice := strings.SplitAfter(sslCertContent, "-----END CERTIFICATE-----\n")
	for _, sslCertContentStr := range sslCertContentSlice {
		certificate, err := entity.NewSslPair(sslCertContentStr)
		if err != nil {
			return certificates, errors.New("SslCertificateError")
		}
		certificates = append(certificates, certificate)
	}

	return certificates, nil
}

func (repo SslQueryRepo) SslFactory(
	sslHostname string,
	sslPrivateKey string,
	sslCertContent string,
) (entity.Ssl, error) {
	var ssl entity.Ssl

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

	var certificate entity.SslPair
	var chainCertificates []entity.SslPair
	for _, sslCertificate := range sslCertificates {
		if sslCertificate.IsCA {
			certificate = sslCertificate
			continue
		}

		chainCertificates = append(chainCertificates, sslCertificate)
	}

	id, err := valueObject.NewSslId(certificate.SerialNumber.String())
	if err != nil {
		return ssl, err
	}

	return entity.NewSsl(
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
		"sed", "-n", "/virtualhost/, /}/p", repo.olsHttpdConfigPath,
	)
	if err != nil {
		matchErr, _ := regexp.MatchString("No such file or directory", err.Error())
		if !matchErr {
			return []HttpdVhostConfig{}, err
		}

		return []HttpdVhostConfig{}, errors.New("HttpdVhostsConfigEmpty")
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

func (repo SslQueryRepo) Get() ([]entity.Ssl, error) {
	var ssls []entity.Ssl
	httpdVhostsConfig, err := repo.GetHttpdVhostsConfig()
	if err != nil {
		return []entity.Ssl{}, err
	}

	for _, httpdVhostConfig := range httpdVhostsConfig {
		vhostConfigOutput, err := infraHelper.RunCmd(
			"sed", "-n", "/vhssl/, /}/p", httpdVhostConfig.FilePath,
		)
		if err != nil {
			return []entity.Ssl{}, err
		}

		if len(vhostConfigOutput) < 1 {
			return []entity.Ssl{}, nil
		}

		vhostConfigGroups := infraHelper.GetRegexNamedGroups(vhostConfigOutput, "(?:keyFile\\s*(?P<keyFile>.*))?\n\\s*(?:certFile\\s*(?P<certFile>.*))\n\\s*(?:certChain\\s*(?P<certChain>.*))\n\\s*(?:(?:CACertPath\\s*(?P<CACertPath>.*))\n\\s*)?(?:(?:CACertFile\\s*(?P<CACertFile>.*))?)")
		privateKeyOutput, err := infraHelper.RunCmd("cat", vhostConfigGroups["keyFile"])
		if err != nil {
			return []entity.Ssl{}, err
		}

		certFileOutput, err := infraHelper.RunCmd("cat", vhostConfigGroups["certFile"])
		if err != nil {
			return []entity.Ssl{}, err
		}

		ssl, err := repo.SslFactory(httpdVhostConfig.VirtualHost, privateKeyOutput, certFileOutput)
		if err != nil {
			return []entity.Ssl{}, err
		}

		ssls = append(ssls, ssl)
	}

	return ssls, nil
}

func (repo SslQueryRepo) GetById(sslId valueObject.SslId) (entity.Ssl, error) {
	ssls, err := repo.Get()
	if err != nil {
		return entity.Ssl{}, err
	}

	if len(ssls) < 1 {
		return entity.Ssl{}, errors.New("SslNotFound")
	}

	for _, ssl := range ssls {
		if ssl.Id.String() != sslId.String() {
			continue
		}

		return ssl, nil
	}

	return entity.Ssl{}, errors.New("SslNotFound")
}
