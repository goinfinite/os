package infra

import (
	"errors"
	"os"
	"regexp"
	"strconv"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

type SslCmdRepo struct{}

func (repo SslCmdRepo) Add(addSsl dto.AddSsl) error {
	sslQueryRepo := NewSslQueryRepo()

	httpdVhostsConfig, err := sslQueryRepo.GetHttpdVhostsConfig()
	if err != nil {
		matchErr, _ := regexp.MatchString("^(HttpdVhostsConfigEmpty|VhostConfigEmpty)$", err.Error())
		if !matchErr {
			return err
		}

		return errors.New("HttpdVhostsConfigEmpty")
	}

	defaultAddSslId := 1
	sslId, err := valueObject.NewSslId(defaultAddSslId)
	if err != nil {
		return err
	}

	sslCertDirPath := "/speedia/pki/" + addSsl.Hostname.String()
	sslCertFilePath := sslCertDirPath + "/ssl.crt"
	sslPrivateKeyFilePath := sslCertDirPath + "/ssl.key"

	err = infraHelper.MakeDir(sslCertDirPath)
	if err != nil {
		return err
	}

	err = infraHelper.UpdateFile(sslCertFilePath, addSsl.Certificate.String(), true)
	if err != nil {
		return err
	}

	err = infraHelper.UpdateFile(sslPrivateKeyFilePath, addSsl.Key.String(), true)
	if err != nil {
		return err
	}

	newSsl, err := sslQueryRepo.SslFactory(
		int(sslId.Get()),
		addSsl.Hostname.String(),
		addSsl.Key.String(),
		addSsl.Certificate.String(),
	)
	if err != nil {
		return err
	}

	isChainedCert := 0
	if len(newSsl.ChainCertificates) > 1 {
		isChainedCert = 1
	}

	caCertPath := ""
	if isChainedCert == 1 {
		caCertPath = "\n\tCACertPath\t" + sslCertFilePath
	}

	caCertFile := ""
	if isChainedCert == 1 {
		caCertPath = "\n\tCACertFile\t" + sslCertFilePath
	}

	for _, httpdVhostConfig := range httpdVhostsConfig {
		if httpdVhostConfig.VirtualHost != addSsl.Hostname.String() {
			continue
		}

		err = infraHelper.UpdateFile(
			httpdVhostConfig.FilePath,
			"\n\nvhssl {\n\tkeyFile\t"+sslPrivateKeyFilePath+"\n\tcertFile\t"+sslCertFilePath+"\n\tcertChain\t"+strconv.Itoa(isChainedCert)+caCertPath+caCertFile+"\n}",
			false,
		)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (repo SslCmdRepo) Delete(sslId valueObject.SslId) error {
	sslQueryRepo := NewSslQueryRepo()

	sslToDelete, err := sslQueryRepo.GetById(sslId)
	if err != nil {
		return errors.New("SslNotFound")
	}

	httpdVhostsConfig, err := sslQueryRepo.GetHttpdVhostsConfig()
	if err != nil {
		matchErr, _ := regexp.MatchString("^(HttpdVhostsConfigEmpty|VhostConfigEmpty)$", err.Error())
		if !matchErr {
			return err
		}

		return errors.New("HttpdVhostsConfigEmpty")
	}

	sslCertDirPath := "/speedia/pki/" + sslToDelete.Hostname.String()
	err = os.RemoveAll(sslCertDirPath)
	if err != nil {
		return err
	}

	for _, httpdVhostConfig := range httpdVhostsConfig {
		if httpdVhostConfig.VirtualHost != sslToDelete.Hostname.String() {
			continue
		}

		vhostConfigOutput, err := infraHelper.RunCmd(
			"cat",
			httpdVhostConfig.FilePath,
		)

		matchVhostConfigVhssl := regexp.MustCompile(`vhssl\s*\{[^}]*\}`)
		vhostConfigWithoutVhssl := matchVhostConfigVhssl.ReplaceAllString(vhostConfigOutput, "")

		err = infraHelper.UpdateFile(httpdVhostConfig.FilePath, vhostConfigWithoutVhssl, true)
		if err != nil {
			return err
		}
		break
	}
	return nil
}
