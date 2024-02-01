package sslInfra

import (
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type SslCmdRepo struct{}

func (repo SslCmdRepo) Add(addSslPair dto.AddSslPair) error {
	sslQueryRepo := SslQueryRepo{}

	vhostThatHasSslHardLink := addSslPair.VirtualHosts[0]
	sslPair, err := sslQueryRepo.GetSslPairByVirtualHost(vhostThatHasSslHardLink)
	if err == nil {
		err = repo.Delete(sslPair.Id)
		if err != nil {
			return err
		}
	}

	vhostThatHasSslHardLinkStr := vhostThatHasSslHardLink.String()
	shouldOverwrite := true

	vhostCertFilePath := "/app/conf/pki/" + vhostThatHasSslHardLinkStr + ".crt"
	err = infraHelper.UpdateFile(vhostCertFilePath, addSslPair.Certificate.String(), shouldOverwrite)
	if err != nil {
		return err
	}

	vhostCertKeyFilePath := "/app/conf/pki/" + vhostThatHasSslHardLinkStr + ".key"
	err = infraHelper.UpdateFile(vhostCertKeyFilePath, addSslPair.Key.String(), shouldOverwrite)
	if err != nil {
		return err
	}

	isSymlink := false
	for _, vhost := range addSslPair.VirtualHosts {
		if vhost.String() != vhostThatHasSslHardLinkStr {
			isSymlink = true
		}

		if isSymlink {
			vhostCertSymlinkPath := "/app/conf/pki/" + vhostThatHasSslHardLinkStr + ".crt"
			err = os.Symlink(vhostCertFilePath, vhostCertSymlinkPath)
			if err != nil {
				log.Printf("AddSslPairError (%s): %s", vhost.String(), err.Error())
			}

			vhostCertKeySymlinkPath := "/app/conf/pki/" + vhostThatHasSslHardLinkStr + ".key"
			err = os.Symlink(vhostCertKeyFilePath, vhostCertKeySymlinkPath)
			if err != nil {
				log.Printf("AddSslPairError (%s): %s", vhost.String(), err.Error())
			}
		}

		vhostConfFilePath, err := sslQueryRepo.GetVhostConfFilePath(vhost)
		if err != nil {
			log.Printf("AddSslPairError (%s): %s", vhost.String(), err.Error())
		}

		_, err = infraHelper.RunCmd(
			"sed",
			"-i",
			"/root \\/app\\/html\\/"+vhost.String()+";/a\\\\n"+
				"    ssl_certificate /app/conf/pki/"+vhost.String()+".crt;\\n"+
				"    ssl_certificate_key /app/conf/pki/"+vhost.String()+".key;\\n",
			vhostConfFilePath.String(),
		)
		if err != nil {
			log.Printf("AddSslPairError (%s): %s", vhost.String(), err.Error())
		}
	}

	return nil
}

func (repo SslCmdRepo) deleteSslConfByVhost(vhost valueObject.Fqdn) error {
	sslQueryRepo := SslQueryRepo{}
	vhostConfFilePath, err := sslQueryRepo.GetVhostConfFilePath(vhost)
	if err != nil {
		return err
	}

	vhostConfContentStr, err := infraHelper.GetFileContent(vhostConfFilePath.String())
	if err != nil {
		return err
	}

	vhostSslPortConfRegex := regexp.MustCompile(`\s*listen 443 ssl;`)
	vhostConfWithoutSslPort := vhostSslPortConfRegex.ReplaceAllString(vhostConfContentStr, "")
	vhostSslConfRegex := regexp.MustCompile(
		`\s*ssl_certificate\s+[^\n]*\n\s*ssl_certificate_key\s+[^\n]*\n`,
	)
	vhostConfWithoutSslConf := vhostSslConfRegex.ReplaceAllString(vhostConfWithoutSslPort, "")
	return infraHelper.UpdateFile(vhostConfFilePath.String(), vhostConfWithoutSslConf, true)
}

func (repo SslCmdRepo) Delete(sslId valueObject.SslId) error {
	sslQueryRepo := SslQueryRepo{}

	sslPairsToDelete, err := sslQueryRepo.GetSslPairById(sslId)
	if err != nil {
		return errors.New("SslNotFound")
	}

	vhostThatHasSslHardLink := sslPairsToDelete.VirtualHosts[0]
	err = repo.deleteSslConfByVhost(vhostThatHasSslHardLink)
	if err != nil {
		return err
	}

	vhostThatHasSslHardLinkStr := vhostThatHasSslHardLink.String()

	vhostCertFilePath := "/app/conf/pki/" + vhostThatHasSslHardLinkStr + ".crt"
	err = os.Remove(vhostCertFilePath)
	if err != nil {
		return err
	}

	vhostCertKeyFilePath := "/app/conf/pki/" + vhostThatHasSslHardLinkStr + ".key"
	err = os.Remove(vhostCertKeyFilePath)
	if err != nil {
		return err
	}

	if len(sslPairsToDelete.VirtualHosts) == 1 {
		return nil
	}

	for _, sslPairVhostToDelete := range sslPairsToDelete.VirtualHosts {
		err = repo.deleteSslConfByVhost(vhostThatHasSslHardLink)
		if err != nil {
			log.Printf("DeleteSslError (%s): %s", sslPairVhostToDelete.String(), err.Error())
			continue
		}

		vhostCertSymlinkPath := "/app/conf/pki/" + vhostThatHasSslHardLinkStr + ".crt"
		err = os.Remove(vhostCertSymlinkPath)
		if err != nil {
			log.Printf("DeleteSslError (%s): %s", sslPairVhostToDelete.String(), err.Error())
			continue
		}

		vhostCertKeySymlinkPath := "/app/conf/pki/" + vhostThatHasSslHardLinkStr + ".key"
		err = os.Remove(vhostCertKeySymlinkPath)
		if err != nil {
			log.Printf("DeleteSslError (%s): %s", sslPairVhostToDelete.String(), err.Error())
			continue
		}

		log.Printf(
			"SSL '%s' of '%s' virtual host deleted.",
			sslId.String(),
			sslPairVhostToDelete.String(),
		)
	}

	return nil
}
