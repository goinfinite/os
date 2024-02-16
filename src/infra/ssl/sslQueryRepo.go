package sslInfra

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

type SslQueryRepo struct{}

type SslCertificates struct {
	MainCertificate     entity.SslCertificate
	ChainedCertificates []entity.SslCertificate
}

func (repo SslQueryRepo) sslCertificatesFactory(
	sslCertContent valueObject.SslCertificateContent,
) (SslCertificates, error) {
	var certificates SslCertificates

	sslCertContentSlice := strings.SplitAfter(
		sslCertContent.String(),
		"-----END CERTIFICATE-----\n",
	)
	for _, sslCertContentStr := range sslCertContentSlice {
		if len(sslCertContentStr) == 0 {
			continue
		}

		certificateContent, err := valueObject.NewSslCertificateContent(sslCertContentStr)
		if err != nil {
			return certificates, err
		}

		certificate, err := entity.NewSslCertificate(certificateContent)
		if err != nil {
			return certificates, err
		}

		if !certificate.IsCA {
			certificates.MainCertificate = certificate
			continue
		}

		certificates.ChainedCertificates = append(certificates.ChainedCertificates, certificate)
	}

	return certificates, nil
}

func (repo SslQueryRepo) sslPairFactory(
	sslVhosts []valueObject.Fqdn,
) (entity.SslPair, error) {
	var ssl entity.SslPair

	firstVhost := sslVhosts[0]
	firstVhostStr := firstVhost.String()

	vhostCertKeyFilePath := pkiConfDir + "/" + firstVhostStr + ".key"
	vhostCertKeyContentStr, err := infraHelper.GetFileContent(vhostCertKeyFilePath)
	if err != nil {
		return ssl, errors.New("FailedToOpenCertKeyFile: " + err.Error())
	}
	privateKey, err := valueObject.NewSslPrivateKey(vhostCertKeyContentStr)
	if err != nil {
		return ssl, err
	}

	vhostCertFilePath := pkiConfDir + "/" + firstVhostStr + ".crt"
	vhostCertFileContentStr, err := infraHelper.GetFileContent(vhostCertFilePath)
	if err != nil {
		return ssl, errors.New("FailedToOpenCertFile: " + err.Error())
	}
	certificate, err := valueObject.NewSslCertificateContent(vhostCertFileContentStr)
	if err != nil {
		return ssl, err
	}

	sslCertificates, err := repo.sslCertificatesFactory(certificate)
	if err != nil {
		return ssl, errors.New("FailedToGetMainAndChainedCerts: " + err.Error())
	}

	mainCertificate := sslCertificates.MainCertificate
	chainCertificates := sslCertificates.ChainedCertificates

	var chainCertificatesContent []valueObject.SslCertificateContent
	for _, sslChainCertificate := range chainCertificates {
		chainCertificatesContent = append(
			chainCertificatesContent,
			sslChainCertificate.CertificateContent,
		)
	}

	hashId, err := valueObject.NewSslIdFromSslPairContent(
		mainCertificate.CertificateContent,
		chainCertificatesContent,
		privateKey,
	)
	if err != nil {
		return ssl, err
	}

	return entity.NewSslPair(
		hashId,
		sslVhosts,
		mainCertificate,
		privateKey,
		chainCertificates,
	), nil
}

func (repo SslQueryRepo) GetSslPairs() ([]entity.SslPair, error) {
	// Criar um slice vazio de SslPairs
	sslPairs := []entity.SslPair{}

	// Criar um map cujo a chave é o caminho do arquivo ".crt" original (target) e o valor é um slice de FQDN que terá o host original (target) e os symlinks
	certFilePathWithVhosts := map[string][]valueObject.Fqdn{}
	// Buscar todos os vhosts
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhosts, err := vhostQueryRepo.Get()
	if err != nil {
		return sslPairs, errors.New("FailedToGetVhosts")
	}
	// Iterar sobre os vhosts
	for _, vhost := range vhosts {
		// Montar o caminho do arquivo ".crt" utilizando o vhost da iteração atual
		certFilePath := pkiConfDir + "/" + vhost.Hostname.String() + ".crt"
		// Validar se ele é um symlink
		isSymlink := infraHelper.IsSymlink(certFilePath)
		// Se for, descobrir qual o caminho do arquivo ".crt" original (target) do vhost da iteração atual através do os.Readlink()
		if isSymlink {
			targetCertFilePath, err := os.Readlink(certFilePath)
			if err != nil {
				log.Printf("FailedToGetTargetCertFilePathFromSymlink: %s", err.Error())
				continue
			}

			certFilePath = targetCertFilePath
		}
		// Validar se o caminho do arquivo ".crt" já existe no map
		_, certFilePathAlreadyExistsInMap := certFilePathWithVhosts[certFilePath]
		// Se não existir, adicionar como chave  com um slice de FQDN vazio como valor
		if !certFilePathAlreadyExistsInMap {
			certFilePathWithVhosts[certFilePath] = []valueObject.Fqdn{}
		}
		// Adicionar o vhost da iteração atual como valor ao map com a chave igual ao caminho do arquivo ".crt" da iteração atual
		certFilePathWithVhosts[certFilePath] = append(
			certFilePathWithVhosts[certFilePath],
			vhost.Hostname,
		)
	}

	// Iterar sobre o map
	for certFilePath, vhosts := range certFilePathWithVhosts {
		// Enviar o slice de vhosts da iteração atual para o factory
		sslPair, err := repo.sslPairFactory(vhosts)
		if err != nil {
			log.Printf("FailedToGetSslPair (%s): %s", certFilePath, err.Error())
			continue
		}

		// Adicionar o SslPair retornado da factory ao slice de SslPairs
		sslPairs = append(sslPairs, sslPair)
	}

	// Retornar os SslPairs
	return sslPairs, nil
}

func (repo SslQueryRepo) GetSslPairById(sslId valueObject.SslId) (entity.SslPair, error) {
	sslPairs, err := repo.GetSslPairs()
	if err != nil {
		return entity.SslPair{}, err
	}

	if len(sslPairs) < 1 {
		return entity.SslPair{}, errors.New("SslPairNotFound")
	}

	for _, ssl := range sslPairs {
		if ssl.Id.String() != sslId.String() {
			continue
		}

		return ssl, nil
	}

	return entity.SslPair{}, errors.New("SslPairNotFound")
}
