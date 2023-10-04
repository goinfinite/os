package infra

import (
	"errors"
	"log"
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

func NewSslQueryRepo() *SslQueryRepo {
	return &SslQueryRepo{
		olsHttpdConfigPath: "/usr/local/lsws/conf/httpd_config.conf",
	}
}

func (repo SslQueryRepo) SslFactory(
	sslId int,
	sslHostname string,
	sslPrivateKey string,
	sslCertContent string,
) (entity.Ssl, error) {
	var ssl entity.Ssl
	id, err := valueObject.NewSslId(sslId)
	if err != nil {
		return ssl, errors.New("SslIdError")
	}

	privateKey, err := valueObject.NewSslPrivateKey(sslPrivateKey)
	if err != nil {
		return ssl, errors.New("SslPrivateKeyError")
	}

	var certificate valueObject.SslCertificate
	var chainCertificates []valueObject.SslCertificate
	sslCertContentSlice := strings.SplitAfter(sslCertContent, "-----END CERTIFICATE-----\n")
	for sslCertContentIndex, sslCertContentStr := range sslCertContentSlice {
		if sslCertContentIndex == 0 {
			certificate, err = valueObject.NewSslCertificate(sslCertContentStr)
			if err != nil {
				return ssl, errors.New("SslCertificateError")
			}
			continue
		}

		chainCertificate, err := valueObject.NewSslCertificate(sslCertContentStr)
		if err != nil {
			return ssl, errors.New("SslCertificateError")
		}
		chainCertificates = append(chainCertificates, chainCertificate)
	}

	hostname, err := valueObject.NewVirtualHost(sslHostname)
	if err != nil {
		return ssl, errors.New("SslHostnameError")
	}

	certInfo, _ := certificate.GetCertInfo()

	certIssuedAtUnix := certInfo.NotBefore.Unix()
	issuedAt := valueObject.UnixTime(certIssuedAtUnix)

	certExpireAtUnix := certInfo.NotAfter.Unix()
	expireAt := valueObject.UnixTime(certExpireAtUnix)

	return entity.NewSsl(
		id,
		hostname,
		&issuedAt,
		&expireAt,
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

func (repo SslQueryRepo) Get() ([]entity.Ssl, error) {
	var ssls []entity.Ssl
	httpdVhostsConfig, err := repo.GetHttpdVhostsConfig()
	if err != nil {
		return []entity.Ssl{}, err
	}

	for httpdVhostConfigIndex, httpdVhostConfig := range httpdVhostsConfig {
		vhostConfigOutput, err := infraHelper.RunCmd(
			"sed", "-n", "/vhssl/, /}/p", httpdVhostConfig.FilePath,
		)
		if err != nil {
			log.Print(err)
			return []entity.Ssl{}, err
		}

		if len(vhostConfigOutput) < 1 {
			return []entity.Ssl{}, errors.New("VhostConfigEmpty")
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

		sslId := httpdVhostConfigIndex + 1
		ssl, err := repo.SslFactory(sslId, httpdVhostConfig.VirtualHost, privateKeyOutput, certFileOutput)
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
