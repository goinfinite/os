package sslInfra

import (
	"errors"
	"log"
	"os"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

const pkiConfDir = "/app/conf/pki"

type SslCmdRepo struct {
	sslQueryRepo SslQueryRepo
}

func NewSslCmdRepo() SslCmdRepo {
	return SslCmdRepo{
		sslQueryRepo: SslQueryRepo{},
	}
}

func (repo SslCmdRepo) forceSymlink(
	pkiSourcePath string,
	pkiTargetPath string,
) error {
	err := os.Remove(pkiTargetPath)
	if err != nil {
		return errors.New("FailedToDeletePkiFile: " + err.Error())
	}

	err = os.Symlink(pkiSourcePath, pkiTargetPath)
	if err != nil {
		return errors.New("AddPkiSymlinkError: " + err.Error())
	}

	return nil
}

func (repo SslCmdRepo) replaceWithSelfSigned(vhost valueObject.Fqdn) error {
	vhostStr := vhost.String()

	vhostCertFilePath := pkiConfDir + "/" + vhostStr + ".crt"
	err := os.Remove(vhostCertFilePath)
	if err != nil {
		return errors.New("FailedToDeleteCertFile: " + err.Error())
	}

	vhostCertKeyFilePath := pkiConfDir + "/" + vhostStr + ".key"
	err = os.Remove(vhostCertKeyFilePath)
	if err != nil {
		return errors.New("FailedToDeleteCertKeyFile: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"openssl",
		"req",
		"-x509",
		"-nodes",
		"-days",
		"365",
		"-newkey",
		"rsa:2048",
		"-keyout",
		vhostCertKeyFilePath,
		"-out",
		vhostCertFilePath,
		"-subj",
		"/C=US/ST=California/L=LosAngeles/O=Acme/CN="+vhostStr,
	)
	if err != nil {
		return errors.New("ReplaceWithSelfSignedFailed: " + err.Error())
	}

	return nil
}

func (repo SslCmdRepo) Add(addSslPair dto.AddSslPair) error {
	if len(addSslPair.VirtualHosts) == 0 {
		return errors.New("NoVirtualHostsProvidedToAddSslPair")
	}

	firstVhostStr := addSslPair.VirtualHosts[0].String()
	firstVhostCertFilePath := pkiConfDir + "/" + firstVhostStr + ".crt"
	firstVhostCertKeyFilePath := pkiConfDir + "/" + firstVhostStr + ".key"

	for _, vhost := range addSslPair.VirtualHosts {
		vhostStr := vhost.String()
		vhostCertFilePath := pkiConfDir + "/" + vhostStr + ".crt"
		vhostCertKeyFilePath := pkiConfDir + "/" + vhostStr + ".key"

		shouldBeSymlink := vhostStr != firstVhostStr
		if shouldBeSymlink {
			err := repo.forceSymlink(firstVhostCertFilePath, vhostCertFilePath)
			if err != nil {
				log.Printf("AddSslCertSymlinkError (%s): %s", vhost.String(), err.Error())
				continue
			}

			err = repo.forceSymlink(firstVhostCertKeyFilePath, vhostCertKeyFilePath)
			if err != nil {
				log.Printf("AddSslKeySymlinkError (%s): %s", vhost.String(), err.Error())
				continue
			}

			continue
		}

		shouldOverwrite := true
		err := infraHelper.UpdateFile(
			vhostCertFilePath,
			addSslPair.Certificate.CertificateContent.String(),
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
	}

	return nil
}

func (repo SslCmdRepo) Delete(sslId valueObject.SslId) error {
	sslPairToDelete, err := repo.sslQueryRepo.GetSslPairById(sslId)
	if err != nil {
		return errors.New("SslNotFound")
	}

	for _, vhost := range sslPairToDelete.VirtualHosts {
		err = repo.replaceWithSelfSigned(vhost)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), vhost.String())
			continue
		}
	}

	return nil
}
