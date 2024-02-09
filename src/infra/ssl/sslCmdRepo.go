package sslInfra

import (
	"errors"
	"log"
	"os"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type SslCmdRepo struct {
	sslQueryRepo SslQueryRepo
}

func NewSslCmdRepo() SslCmdRepo {
	return SslCmdRepo{
		sslQueryRepo: SslQueryRepo{},
	}
}

func (repo SslCmdRepo) SelfSignedCertSwap(vhost valueObject.Fqdn) error {
	selfSignedSslKeyPath := "/app/conf/pki/" + vhost.String() + ".key"
	selfSignedSslCertPath := "/app/conf/pki/" + vhost.String() + ".crt"

	_, err := infraHelper.RunCmd(
		"openssl",
		"req",
		"-x509",
		"-nodes",
		"-days",
		"365",
		"-newkey",
		"rsa:2048",
		"-keyout",
		selfSignedSslKeyPath,
		"-out",
		selfSignedSslCertPath,
		"-subj",
		"/C=US/ST=California/L=LosAngeles/O=Acme/CN="+vhost.String(),
	)
	if err != nil {
		return errors.New("SelfSignedCertSwapFailed: " + err.Error())
	}

	return nil
}

func (repo SslCmdRepo) Add(addSslPair dto.AddSslPair) error {
	if len(addSslPair.VirtualHosts) == 0 {
		return errors.New("NoVirtualHostsProvidedToAddSslPair")
	}

	firstVhostStr := addSslPair.VirtualHosts[0].String()
	for _, vhost := range addSslPair.VirtualHosts {
		vhostStr := vhost.String()
		vhostCertFilePath := "/app/conf/pki/" + vhostStr + ".crt"
		vhostCertKeyFilePath := "/app/conf/pki/" + vhostStr + ".key"

		shouldBeSymlink := vhostStr != firstVhostStr
		if shouldBeSymlink {
			firstVhostCertFilePath := "/app/conf/pki/" + firstVhostStr + ".crt"
			firstVhostCertKeyFilePath := "/app/conf/pki/" + firstVhostStr + ".key"

			err := os.Symlink(firstVhostCertFilePath, vhostCertFilePath)
			if err != nil {
				log.Printf("AddSslCertSymlinkError (%s): %s", vhost.String(), err.Error())
				continue
			}

			err = os.Symlink(firstVhostCertKeyFilePath, vhostCertKeyFilePath)
			if err != nil {
				log.Printf("AddSslKeySymlinkError (%s): %s", vhost.String(), err.Error())
				continue
			}

			continue
		}

		shouldOverwrite := true
		err := infraHelper.UpdateFile(
			vhostCertFilePath,
			addSslPair.Certificate.String(),
			shouldOverwrite,
		)
		if err != nil {
			return err
		}

		err = infraHelper.UpdateFile(
			vhostCertKeyFilePath,
			addSslPair.Key.String(),
			shouldOverwrite,
		)
		if err != nil {
			return err
		}

		log.Printf(
			"SSL '%s' created in '%s' virtual host.",
			addSslPair.Certificate.Id.String(),
			vhost.String(),
		)
	}

	return nil
}

func (repo SslCmdRepo) Delete(sslId valueObject.SslId) error {
	sslPairToDelete, err := repo.sslQueryRepo.GetSslPairById(sslId)
	if err != nil {
		return errors.New("SslNotFound")
	}

	for _, vhost := range sslPairToDelete.VirtualHosts {
		vhostStr := vhost.String()

		vhostCertFilePath := "/app/conf/pki/" + vhostStr + ".crt"
		err = os.Remove(vhostCertFilePath)
		if err != nil {
			log.Printf(
				"FailedToDeleteCertFile (%s): %s", vhostStr, err.Error(),
			)
			continue
		}

		vhostCertKeyFilePath := "/app/conf/pki/" + vhostStr + ".key"
		err = os.Remove(vhostCertKeyFilePath)
		if err != nil {
			log.Printf(
				"FailedToDeleteCertKeyFile (%s): %s", vhostStr, err.Error(),
			)
			continue
		}

		err = repo.SelfSignedCertSwap(vhost)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), vhostStr)
			continue
		}

		sslPair, err := repo.sslQueryRepo.GetSslPairByVirtualHost(vhost)
		if err != nil {
			log.Printf("FailedToGetSelfSignedSsl (%s): %s", vhostStr, err.Error())
			continue
		}

		log.Printf(
			"Self Signed SSL '%s' created in '%s' virtual host.",
			sslPair.Id.String(),
			vhostStr,
		)
	}

	return nil
}
