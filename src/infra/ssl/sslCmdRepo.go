package sslInfra

import (
	"errors"
	"log"
	"os"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type SslCmdRepo struct{}

func (repo SslCmdRepo) GenerateSelfSignedCert(vhost valueObject.Fqdn) error {
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
		return errors.New("GenerateSelfSignedCertFailed: " + err.Error())
	}

	return nil
}

func (repo SslCmdRepo) Add(addSslPair dto.AddSslPair) error {
	sslQueryRepo := SslQueryRepo{}

	if len(addSslPair.VirtualHosts) == 0 {
		return errors.New("NoVirtualHostsProvidedToAddSslPair")
	}

	firstVhostStr := addSslPair.VirtualHosts[0].String()
	for _, vhost := range addSslPair.VirtualHosts {
		sslPair, err := sslQueryRepo.GetSslPairByVirtualHost(vhost)
		if err != nil && err.Error() != "SslPairNotFound" {
			log.Printf("FailedToValidateSslPairExistence (%s): %s", vhost.String(), err.Error())
			continue
		}

		sslPairExists := sslPair.Id != ""
		if sslPairExists {
			err := repo.Delete(sslPair.Id)
			if err != nil {
				log.Printf("FailedToDeleteOldSslPair (%s): %s", vhost.String(), err.Error())
				continue
			}
		}

		vhostStr := vhost.String()
		vhostCertFilePath := "/app/conf/pki/" + vhostStr + ".crt"
		vhostCertKeyFilePath := "/app/conf/pki/" + vhostStr + ".key"

		shouldBeSymlink := vhostStr != firstVhostStr
		if shouldBeSymlink {
			firstVhostCertFilePath := "/app/conf/pki/" + firstVhostStr + ".crt"
			firstVhostCertKeyFilePath := "/app/conf/pki/" + firstVhostStr + ".key"

			err = os.Symlink(firstVhostCertFilePath, vhostCertFilePath)
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
		err = infraHelper.UpdateFile(
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
			"SSL '%v' added in '%v' virtual host.",
			addSslPair.Certificate.Id.String(),
			vhost.String(),
		)
	}

	return nil
}

func (repo SslCmdRepo) Delete(sslId valueObject.SslId) error {
	sslQueryRepo := SslQueryRepo{}
	sslPairToDelete, err := sslQueryRepo.GetSslPairById(sslId)
	if err != nil {
		return errors.New("SslNotFound")
	}

	for _, sslPairVhostToDelete := range sslPairToDelete.VirtualHosts {
		sslPairVhostToDeleteStr := sslPairVhostToDelete.String()

		vhostCertFilePath := "/app/conf/pki/" + sslPairVhostToDeleteStr + ".crt"
		err = os.Remove(vhostCertFilePath)
		if err != nil {
			log.Printf(
				"FailedToDeleteCertFile (%s): %s", sslPairVhostToDeleteStr, err.Error(),
			)
			continue
		}

		vhostCertKeyFilePath := "/app/conf/pki/" + sslPairVhostToDeleteStr + ".key"
		err = os.Remove(vhostCertKeyFilePath)
		if err != nil {
			log.Printf(
				"FailedToDeleteCertKeyFile (%s): %s", sslPairVhostToDeleteStr, err.Error(),
			)
			continue
		}

		err = repo.GenerateSelfSignedCert(sslPairVhostToDelete)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), sslPairVhostToDeleteStr)
			continue
		}

		log.Printf(
			"SSL '%s' of '%s' virtual host deleted.",
			sslId.String(),
			sslPairVhostToDeleteStr,
		)
	}

	return nil
}
