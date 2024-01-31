package sslInfra

import (
	"errors"
	"os"
	"regexp"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type SslCmdRepo struct{}

func (repo SslCmdRepo) Add(addSslPair dto.AddSslPair) error {
	sslQueryRepo := SslQueryRepo{}

	sslPair, err := sslQueryRepo.GetSslPairByVirtualHost(addSslPair.VirtualHost)
	if err == nil {
		err = repo.Delete(sslPair.Id)
		if err != nil {
			return err
		}
	}

	vhostConfFilePath, err := sslQueryRepo.GetVhostConfFilePath(addSslPair.VirtualHost)
	if err != nil {
		return err
	}

	vhostStr := addSslPair.VirtualHost.String()

	sslCertFilePath := "/app/conf/pki/" + vhostStr + ".crt"
	err = infraHelper.UpdateFile(sslCertFilePath, addSslPair.Certificate.String(), true)
	if err != nil {
		return err
	}

	sslKeyFilePath := "/app/conf/pki/" + vhostStr + ".key"
	err = infraHelper.UpdateFile(sslKeyFilePath, addSslPair.Key.String(), true)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"/root \\/app\\/html\\/"+vhostStr+";/a\\\\n"+
			"    ssl_certificate /app/conf/pki/"+vhostStr+".crt;\\n"+
			"    ssl_certificate_key /app/conf/pki/"+vhostStr+".key;\\n",
		vhostConfFilePath.String(),
	)
	if err != nil {
		return err
	}

	return nil
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

	vhostStr := sslToDelete.VirtualHost.String()

	vhostCertFilePath := "/app/conf/pki/" + vhostStr + ".crt"
	err = os.RemoveAll(vhostCertFilePath)
	if err != nil {
		return err
	}

	vhostCertKeyFilePath := "/app/conf/pki/" + vhostStr + ".key"
	err = os.RemoveAll(vhostCertKeyFilePath)
	if err != nil {
		return err
	}

	vhostSslConfRegex := regexp.MustCompile(
		`\s*ssl_certificate\s+[^\n]*\n\s*ssl_certificate_key\s+[^\n]*\n`,
	)
	vhostConfWithoutSsl := vhostSslConfRegex.ReplaceAllString(vhostConfContentStr, "")
	return infraHelper.UpdateFile(vhostConfFilePath.String(), vhostConfWithoutSsl, true)
}
