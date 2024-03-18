package sslInfra

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

const PkiConfDir = "/app/conf/pki"

type SslCmdRepo struct {
	sslQueryRepo SslQueryRepo
}

func NewSslCmdRepo() SslCmdRepo {
	return SslCmdRepo{
		sslQueryRepo: SslQueryRepo{},
	}
}

func (repo SslCmdRepo) deleteCurrentSsl(vhost valueObject.Fqdn) error {
	vhostStr := vhost.String()

	vhostCertFilePath := PkiConfDir + "/" + vhostStr + ".crt"
	vhostCertFileExists := infraHelper.FileExists(vhostCertFilePath)
	if vhostCertFileExists {
		err := os.Remove(vhostCertFilePath)
		if err != nil {
			return errors.New("FailedToDeleteCertFile: " + err.Error())
		}
	}

	vhostCertKeyFilePath := PkiConfDir + "/" + vhostStr + ".key"
	vhostCertKeyFileExists := infraHelper.FileExists(vhostCertKeyFilePath)
	if vhostCertKeyFileExists {
		err := os.Remove(vhostCertKeyFilePath)
		if err != nil {
			return errors.New("FailedToDeleteCertKeyFile: " + err.Error())
		}
	}

	return nil
}

func (repo SslCmdRepo) ReplaceWithSelfSigned(vhost valueObject.Fqdn) error {
	err := repo.deleteCurrentSsl(vhost)
	if err != nil {
		return err
	}

	return infraHelper.CreateSelfSignedSsl(PkiConfDir, vhost.String())
}

func (repo SslCmdRepo) ReplaceWithValidSsl(vhost valueObject.Fqdn) error {
	vhostStr := vhost.String()
	vhostRootDir := "/app/html"
	if !infraHelper.IsVirtualHostPrimaryDomain(vhost) {
		vhostRootDir += "/" + vhostStr
	}

	_, err := infraHelper.RunCmd(
		"certbot",
		"certonly",
		"--webroot",
		"--webroot-path",
		vhostRootDir,
		"--agree-tos",
		"--register-unsafely-without-email",
		"--cert-name",
		vhostStr,
		"-d",
		vhostStr,
	)
	if err != nil {
		return fmt.Errorf("CreateValidSslFailed (%s): %s", vhostStr, err.Error())
	}

	certbotDirPath := "/etc/letsencrypt/live"
	shouldOverwrite := true

	certbotCrtFilePath := certbotDirPath + "/" + vhostStr + "/fullchain.pem"
	vhostCrtFilePath := PkiConfDir + "/" + vhostStr + ".crt"
	err = infraHelper.CreateSymlink(
		certbotCrtFilePath,
		vhostCrtFilePath,
		shouldOverwrite,
	)
	if err != nil {
		return fmt.Errorf("CreateSslCrtSymlinkError (%s): %s", vhost.String(), err.Error())
	}

	certbotKeyFilePath := certbotDirPath + "/" + vhostStr + "/privkey.pem"
	vhostCrtKeyFilePath := PkiConfDir + "/" + vhostStr + ".key"
	err = infraHelper.CreateSymlink(
		certbotKeyFilePath,
		vhostCrtKeyFilePath,
		shouldOverwrite,
	)
	if err != nil {
		return fmt.Errorf("CreateSslKeySymlinkError (%s): %s", vhost.String(), err.Error())
	}

	return nil
}

func (repo SslCmdRepo) Create(createSslPair dto.CreateSslPair) error {
	if len(createSslPair.VirtualHosts) == 0 {
		return errors.New("NoVirtualHostsProvidedToCreateSslPair")
	}

	firstVhostStr := createSslPair.VirtualHosts[0].String()
	firstVhostCertFilePath := PkiConfDir + "/" + firstVhostStr + ".crt"
	firstVhostCertKeyFilePath := PkiConfDir + "/" + firstVhostStr + ".key"

	for _, vhost := range createSslPair.VirtualHosts {
		vhostStr := vhost.String()
		vhostCertFilePath := PkiConfDir + "/" + vhostStr + ".crt"
		vhostCertKeyFilePath := PkiConfDir + "/" + vhostStr + ".key"

		shouldBeSymlink := vhostStr != firstVhostStr
		if shouldBeSymlink {
			shouldOverwrite := true
			err := infraHelper.CreateSymlink(
				firstVhostCertFilePath,
				vhostCertFilePath,
				shouldOverwrite,
			)
			if err != nil {
				log.Printf("CreateSslCertSymlinkError (%s): %s", vhost.String(), err.Error())
				continue
			}

			err = infraHelper.CreateSymlink(
				firstVhostCertKeyFilePath,
				vhostCertKeyFilePath,
				shouldOverwrite,
			)
			if err != nil {
				log.Printf("CreateSslKeySymlinkError (%s): %s", vhost.String(), err.Error())
				continue
			}

			continue
		}

		shouldOverwrite := true
		err := infraHelper.UpdateFile(
			vhostCertFilePath,
			createSslPair.Certificate.CertificateContent.String(),
			shouldOverwrite,
		)
		if err != nil {
			return err
		}

		err = infraHelper.UpdateFile(
			vhostCertKeyFilePath,
			createSslPair.Key.String(),
			shouldOverwrite,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo SslCmdRepo) Delete(sslId valueObject.SslId) error {
	sslPairToDelete, err := repo.sslQueryRepo.GetSslPairById(sslId)
	if err != nil {
		return errors.New("SslNotFound")
	}

	for _, vhost := range sslPairToDelete.VirtualHosts {
		err = repo.ReplaceWithSelfSigned(vhost)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), vhost.String())
			continue
		}
	}

	return nil
}
