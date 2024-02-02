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
	// Laço em cima do addSslPair.VirtualHosts
	// Para cada vhost:
	vhostSymlinkOf := addSslPair.VirtualHosts[0]
	for vhostIndex, vhost := range addSslPair.VirtualHosts {
		// Verificar se o SSL já existe.
		sslPair, err := sslQueryRepo.GetSslPairByVirtualHost(vhost)
		// Caso dê erro que não seja "NotFound", logar e dar continue.
		if err != nil && err.Error() != "SslPairNotFound" {
			log.Printf("FailedToValidateSslPairExistence (%s): %s", vhost.String(), err.Error())
			continue
		}
		// Se existir, chamar o Delete().
		sslPairExists := sslPair.Id != ""
		if sslPairExists {
			err := repo.Delete(sslPair.Id)
			// Caso dê erro, logar e dar continue.
			if err != nil {
				log.Printf("FailedToDeleteTheOldSslPair (%s): %s", vhost.String(), err.Error())
				continue
			}
		}
		// Verificar se é o primeiro vhost
		isSymlink := vhostIndex != 0
		// Se não for, é um symlink
		if isSymlink {
			// Sendo um symlink, deve-se criar um symlink pro cert
			vhostCertToSymlinkPath := "/app/conf/pki/" + vhostSymlinkOf.String() + ".crt"
			vhostCertSymlinkPath := "/app/conf/pki/" + vhost.String() + ".crt"
			err = os.Symlink(vhostCertToSymlinkPath, vhostCertSymlinkPath)
			// Caso dê erro, logar e dar continue
			if err != nil {
				log.Printf("FailedToAddSslCertSymlink (%s): %s", vhost.String(), err.Error())
				continue
			}

			// Sendo um symlink, deve-se criar um symlink pra key
			vhostKeyToSymlinkPath := "/app/conf/pki/" + vhostSymlinkOf.String() + ".key"
			vhostCertKeySymlinkPath := "/app/conf/pki/" + vhost.String() + ".key"
			err = os.Symlink(vhostKeyToSymlinkPath, vhostCertKeySymlinkPath)
			// Caso dê erro, logar e dar continue
			if err != nil {
				log.Printf("FailedToAddSslKeySymlink (%s): %s", vhost.String(), err.Error())
				continue
			}
		}

		if !isSymlink {
			shouldOverwrite := true
			// Usar o UpdateFile para criar ou atualizar o arquivo .crt sem overwrite.
			vhostCertFilePath := "/app/conf/pki/" + vhost.String() + ".crt"
			err = infraHelper.UpdateFile(vhostCertFilePath, addSslPair.Certificate.String(), shouldOverwrite)
			// Caso dê erro, logar e dar continue.
			if err != nil {
				return err
			}
			// Usar o UpdateFile para criar ou atualizar o arquivo .key sem overwrite.
			vhostCertKeyFilePath := "/app/conf/pki/" + vhost.String() + ".key"
			err = infraHelper.UpdateFile(vhostCertKeyFilePath, addSslPair.Key.String(), shouldOverwrite)
			// Caso dê erro, logar e dar continue.
			if err != nil {
				return err
			}
		}

		// Pegar o caminho do arquivo de configuração NGINX do vhost
		vhostConfFilePath, err := sslQueryRepo.GetVhostConfFilePath(vhost)
		// Caso dê erro, logar e dar continue
		if err != nil {
			log.Printf("FailedToGetVhostConfFilePath (%s): %s", vhost.String(), err.Error())
			continue
		}
		// Usar o sed para atualizar o arquivo de configuração do vhost para adicionar os caminhos do .crt e do .key.
		_, err = infraHelper.RunCmd(
			"sed",
			"-i",
			"/root \\/app\\/html\\/"+vhost.String()+";/a\\\\n"+
				"    ssl_certificate /app/conf/pki/"+vhost.String()+".crt;\\n"+
				"    ssl_certificate_key /app/conf/pki/"+vhost.String()+".key;\\n",
			vhostConfFilePath.String(),
		)
		// Caso dê erro, logar e dar continue.
		if err != nil {
			log.Printf("AddSslPairError (%s): %s", vhost.String(), err.Error())
			continue
		}

		// Logar sucesso se tudo der certo pro vhost.
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
	sslPairsToDelete, err := sslQueryRepo.GetSslPairById(sslId)
	if err != nil {
		return errors.New("SslNotFound")
	}

	for _, sslPairVhostToDelete := range sslPairsToDelete.VirtualHosts {
		sslPairVhostToDeleteStr := sslPairVhostToDelete.String()

		vhostCertFilePath := "/app/conf/pki/" + sslPairVhostToDeleteStr + ".crt"
		err = os.Remove(vhostCertFilePath)
		if err != nil {
			log.Printf(
				"FailedToDeleteCertFile (%s): %s", sslPairVhostToDelete.String(), err.Error(),
			)
			continue
		}

		vhostCertKeyFilePath := "/app/conf/pki/" + sslPairVhostToDeleteStr + ".key"
		err = os.Remove(vhostCertKeyFilePath)
		if err != nil {
			log.Printf(
				"FailedToDeleteCertKeyFile (%s): %s", sslPairVhostToDelete.String(), err.Error(),
			)
			continue
		}

		vhostConfFilePath, err := sslQueryRepo.GetVhostConfFilePath(sslPairVhostToDelete)
		if err != nil {
			log.Printf("DeleteSslError (%s): %s", sslPairVhostToDelete.String(), err.Error())
			continue
		}

		vhostConfContentStr, err := infraHelper.GetFileContent(vhostConfFilePath.String())
		if err != nil {
			log.Printf("DeleteSslError (%s): %s", sslPairVhostToDelete.String(), err.Error())
			continue
		}

		vhostSslPortConfRegex := regexp.MustCompile(`\s*listen 443 ssl;`)
		vhostConfWithoutSslPort := vhostSslPortConfRegex.ReplaceAllString(vhostConfContentStr, "")
		vhostSslConfRegex := regexp.MustCompile(
			`\s*ssl_certificate\s+[^\n]*\n\s*ssl_certificate_key\s+[^\n]*\n`,
		)
		vhostConfWithoutSslConf := vhostSslConfRegex.ReplaceAllString(vhostConfWithoutSslPort, "")
		err = infraHelper.UpdateFile(vhostConfFilePath.String(), vhostConfWithoutSslConf, true)
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
