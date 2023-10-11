package infra

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

type SslCmdRepo struct{}

func (repo SslCmdRepo) vhsslConfigFactory(
	sslCertFilePath string,
	sslKeyFilePath string,
	sslHostname string,
	isChained bool,
) (string, error) {
	vhsslChainedConfig := ""
	sslCertChain := "0"
	if isChained {
		sslCertChain = "1"
		vhsslChainedConfig = `
  CACertPath ` + sslCertFilePath + `
  CACertFile ` + sslCertFilePath + ``
	}

	vhsslConfigBreakline := "\n\n"
	vhsslConfig := `
vhssl  {
  keyFile    ` + sslKeyFilePath + `
  certFile   ` + sslCertFilePath + `
  certChain  ` + sslCertChain +
		vhsslChainedConfig + `
}` + vhsslConfigBreakline

	return vhsslConfig, nil
}

func (repo SslCmdRepo) Add(addSsl dto.AddSsl) error {
	sslQueryRepo := SslQueryRepo{}

	vhostConfig, err := sslQueryRepo.GetVhostConfig(addSsl.Hostname.String())
	if err != nil {
		return err
	}

	sslBaseDirPath := "/speedia/pki/" + addSsl.Hostname.String()
	sslKeyFilePath := sslBaseDirPath + "/ssl.key"
	sslCertFilePath := sslBaseDirPath + "/ssl.crt"

	err = infraHelper.MakeDir(sslBaseDirPath)
	if err != nil {
		return err
	}

	sslDtoCert := addSsl.Certificate
	err = infraHelper.UpdateFile(sslCertFilePath, sslDtoCert.Certificate, true)
	if err != nil {
		return err
	}

	sslDtoPk := addSsl.Key
	err = infraHelper.UpdateFile(sslKeyFilePath, sslDtoPk.Key, true)
	if err != nil {
		return err
	}

	newSsl, err := sslQueryRepo.SslFactory(
		addSsl.Hostname.String(),
		sslDtoPk.Key,
		sslDtoCert.Certificate,
	)
	if err != nil {
		return err
	}

	isChainedCert := false
	if len(newSsl.ChainCertificates) > 1 {
		isChainedCert = true
	}

	vhsslConfig, err := repo.vhsslConfigFactory(
		sslCertFilePath,
		sslKeyFilePath,
		addSsl.Hostname.String(),
		isChainedCert,
	)
	err = infraHelper.UpdateFile(vhostConfig.FilePath, vhsslConfig, false)
	if err != nil {
		return err
	}

	return nil
}

func (repo SslCmdRepo) Delete(sslSerialNumber valueObject.SslSerialNumber) error {
	sslQueryRepo := SslQueryRepo{}

	sslToDelete, err := sslQueryRepo.GetById(sslSerialNumber)
	if err != nil {
		return errors.New("SslNotFound")
	}

	vhostConfig, err := sslQueryRepo.GetVhostConfig(sslToDelete.Hostname.String())
	if err != nil {
		return err
	}

	sslBaseDirPath := "/speedia/pki/" + sslToDelete.Hostname.String()
	err = os.RemoveAll(sslBaseDirPath)
	if err != nil {
		return err
	}

	matchVhostConfigVhssl := regexp.MustCompile(`vhssl\s*\{[^}]*\}`)
	vhostConfigWithoutVhssl := matchVhostConfigVhssl.ReplaceAllString(vhostConfig.FileContent, "")

	err = infraHelper.UpdateFile(
		vhostConfig.FilePath,
		strings.TrimRightFunc(vhostConfigWithoutVhssl, unicode.IsSpace)+"\n\n",
		true,
	)
	if err != nil {
		return err
	}

	return nil
}
