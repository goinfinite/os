package sslInfra

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type SslCmdRepo struct{}

func (repo SslCmdRepo) vhsslConfigFactory(
	sslCertFilePath string,
	sslKeyFilePath string,
	isChained bool,
) string {
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
vhssl {
  keyFile    ` + sslKeyFilePath + `
  certFile   ` + sslCertFilePath + `
  certChain  ` + sslCertChain +
		vhsslChainedConfig + `
}` + vhsslConfigBreakline

	return vhsslConfig
}

func (repo SslCmdRepo) Add(addSslPair dto.AddSslPair) error {
	sslQueryRepo := SslQueryRepo{}

	vhostConfFilePath, err := sslQueryRepo.GetVhostConfFilePath(addSslPair.VirtualHost)
	if err != nil {
		return err
	}

	sslBaseDirPath := "/app/conf/pki/" + addSslPair.VirtualHost.String()
	sslKeyFilePath := sslBaseDirPath + "/ssl.key"
	sslCertFilePath := sslBaseDirPath + "/ssl.crt"

	err = infraHelper.MakeDir(sslBaseDirPath)
	if err != nil {
		return err
	}

	err = infraHelper.UpdateFile(sslCertFilePath, addSslPair.Certificate.String(), true)
	if err != nil {
		return err
	}

	err = infraHelper.UpdateFile(sslKeyFilePath, addSslPair.Key.String(), true)
	if err != nil {
		return err
	}

	sslPairCertificate := addSslPair.Certificate
	sslCertificates, err := sslQueryRepo.SslCertificatesFactory(
		sslPairCertificate.Certificate,
	)
	if err != nil {
		return err
	}

	newSsl, err := sslQueryRepo.SslPairFactory(
		addSslPair.VirtualHost,
		addSslPair.Key,
		sslCertificates,
	)
	if err != nil {
		return err
	}

	isChainedCert := true
	if len(newSsl.ChainCertificates) == 1 {
		isChainedCert = false
	}

	vhsslConfig := repo.vhsslConfigFactory(
		sslCertFilePath,
		sslKeyFilePath,
		isChainedCert,
	)
	err = infraHelper.UpdateFile(vhostConfFilePath.String(), vhsslConfig, false)
	return err
}

func (repo SslCmdRepo) Delete(sslId valueObject.SslId) error {
	sslQueryRepo := SslQueryRepo{}

	sslToDelete, err := sslQueryRepo.GetSslPairById(sslId)
	if err != nil {
		return errors.New("SslNotFound")
	}

	vhostConfFilePath, err := sslQueryRepo.GetVhostConfFilePath(sslToDelete.VirtualHost)
	if err != nil {
		return err
	}

	vhostConfContentStr, err := infraHelper.GetFileContent(vhostConfFilePath.String())
	if err != nil {
		return err
	}

	sslBaseDirPath := "/app/conf/pki/" + sslToDelete.VirtualHost.String()
	err = os.RemoveAll(sslBaseDirPath)
	if err != nil {
		return err
	}

	vhostConfVhsslMatch := regexp.MustCompile(`vhssl\s*\{[^}]*\}`)
	vhostConfWithoutVhssl := vhostConfVhsslMatch.ReplaceAllString(vhostConfContentStr, "")
	vhostConfWithoutSpaces := strings.TrimRightFunc(vhostConfWithoutVhssl, unicode.IsSpace)
	vhostConfWithBreakLines := vhostConfWithoutSpaces + "\n\n"

	err = infraHelper.UpdateFile(
		vhostConfFilePath.String(),
		vhostConfWithBreakLines,
		true,
	)
	return err
}
