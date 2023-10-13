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

	vhostConfigFilePath, err := sslQueryRepo.GetVhostConfigFilePath(addSslPair.Hostname)
	if err != nil {
		return err
	}

	sslBaseDirPath := "/speedia/pki/" + addSslPair.Hostname.String()
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

	newSsl, err := sslQueryRepo.SslFactory(
		addSslPair.Hostname.String(),
		addSslPair.Key.String(),
		addSslPair.Certificate.String(),
	)
	if err != nil {
		return err
	}

	isChainedCert := false
	if len(newSsl.ChainCertificates) > 1 {
		isChainedCert = true
	}

	vhsslConfig := repo.vhsslConfigFactory(
		sslCertFilePath,
		sslKeyFilePath,
		isChainedCert,
	)
	err = infraHelper.UpdateFile(vhostConfigFilePath, vhsslConfig, false)
	return err
}

func (repo SslCmdRepo) Delete(sslSerialNumber valueObject.SslSerialNumber) error {
	sslQueryRepo := SslQueryRepo{}

	sslToDelete, err := sslQueryRepo.GetSslPairBySerialNumber(sslSerialNumber)
	if err != nil {
		return errors.New("SslNotFound")
	}

	vhostConfigFilePath, err := sslQueryRepo.GetVhostConfigFilePath(sslToDelete.Hostname)
	if err != nil {
		return err
	}

	vhostConfigBytesOutput, err := os.ReadFile(vhostConfigFilePath)
	if err != nil {
		return err
	}
	vhostConfigOutputStr := string(vhostConfigBytesOutput)

	sslBaseDirPath := "/speedia/pki/" + sslToDelete.Hostname.String()
	err = os.RemoveAll(sslBaseDirPath)
	if err != nil {
		return err
	}

	matchVhostConfigVhssl := regexp.MustCompile(`vhssl\s*\{[^}]*\}`)
	vhostConfigWithoutVhssl := matchVhostConfigVhssl.ReplaceAllString(vhostConfigOutputStr, "")

	err = infraHelper.UpdateFile(
		vhostConfigFilePath,
		strings.TrimRightFunc(vhostConfigWithoutVhssl, unicode.IsSpace)+"\n\n",
		true,
	)
	return err
}
