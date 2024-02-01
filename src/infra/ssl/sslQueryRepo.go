package sslInfra

import (
	"errors"
	"log"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

const configurationsDir = "/app/conf/nginx"

type SslQueryRepo struct{}

type SslCertificates struct {
	MainCertificate     entity.SslCertificate
	ChainedCertificates []entity.SslCertificate
}

func (repo SslQueryRepo) GetVhostConfFilePath(
	vhost valueObject.Fqdn,
) (valueObject.UnixFilePath, error) {
	var vhostConfFilePath valueObject.UnixFilePath

	vhostConfFilePathStr := configurationsDir + "/" + vhost.String() + ".conf"
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	if vhostQueryRepo.IsVirtualHostPrimaryDomain(vhost) {
		vhostConfFilePathStr = configurationsDir + "/primary.conf"
	}

	vhostConfFilePath, err := valueObject.NewUnixFilePath(vhostConfFilePathStr)
	if err != nil {
		return vhostConfFilePath, err
	}

	return vhostConfFilePath, nil
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
	sslPrivateKey valueObject.SslPrivateKey,
	sslCertificates SslCertificates,
) (entity.SslPair, error) {
	var ssl entity.SslPair

	certificate := sslCertificates.MainCertificate
	chainCertificates := sslCertificates.ChainedCertificates

	var chainCertificatesContent []valueObject.SslCertificateContent
	for _, sslChainCertificate := range chainCertificates {
		chainCertificatesContent = append(
			chainCertificatesContent,
			sslChainCertificate.CertificateContent,
		)
	}

	hashId, err := valueObject.NewSslIdFromSslPairContent(
		certificate.CertificateContent,
		chainCertificatesContent,
		sslPrivateKey,
	)
	if err != nil {
		return ssl, err
	}

	return entity.NewSslPair(
		hashId,
		sslVhosts,
		certificate,
		sslPrivateKey,
		chainCertificates,
	), nil
}

func (repo SslQueryRepo) getSymlinkSslPairVhostsByVhost(
	targetVhost valueObject.Fqdn,
) ([]valueObject.Fqdn, error) {
	var vhostsSymlinkedToTargetVhost []valueObject.Fqdn

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	virtualHosts, err := vhostQueryRepo.Get()
	if err != nil {
		return vhostsSymlinkedToTargetVhost, err
	}

	for _, vhost := range virtualHosts {
		isSymlinkOf, err := infraHelper.IsSymlinkTo(
			"/app/conf/pki/"+vhost.Hostname.String()+".crt",
			"/app/conf/pki/"+targetVhost.String()+".crt",
		)
		if err != nil {
			if err.Error() != "FileNotFound" {
				log.Printf("FailedToCheckIfVhostConfIsSymlinkTo (%s): %s", vhost.Hostname.String(), err.Error())
			}

			continue
		}

		if isSymlinkOf {
			vhostsSymlinkedToTargetVhost = append(vhostsSymlinkedToTargetVhost, vhost.Hostname)
		}
	}

	return vhostsSymlinkedToTargetVhost, nil
}

func (repo SslQueryRepo) GetSslPairs() ([]entity.SslPair, error) {
	var sslPairs []entity.SslPair

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	virtualHosts, err := vhostQueryRepo.Get()
	if err != nil {
		return []entity.SslPair{}, err
	}

	for _, vhost := range virtualHosts {
		hostnameStr := vhost.Hostname.String()

		vhostConfigFilePath := configurationsDir + "/" + hostnameStr + ".conf"
		if vhostQueryRepo.IsVirtualHostPrimaryDomain(vhost.Hostname) {
			vhostConfigFilePath = configurationsDir + "/primary.conf"
		}

		vhostConfigContentStr, err := infraHelper.GetFileContent(vhostConfigFilePath)
		if err != nil {
			log.Printf("FailedToOpenVhostConfFile (%s): %s", hostnameStr, err.Error())
			continue
		}

		fileIsEmpty := len(vhostConfigContentStr) < 1
		if fileIsEmpty {
			log.Printf("VirtualHostConfFileIsEmpty (%s)", hostnameStr)
			continue
		}

		vhostCertKeyFilePath := "/app/conf/pki/" + hostnameStr + ".key"
		isSymlink, err := infraHelper.IsSymlink(vhostCertKeyFilePath)
		if err != nil {
			if err.Error() != "FileNotFound" {
				log.Printf(
					"FailedToCheckIfVhostConfIsSymlink (%s): %s",
					vhost.Hostname.String(),
					err.Error(),
				)
			}

			continue
		}

		if isSymlink {
			continue
		}

		sslVhosts := []valueObject.Fqdn{vhost.Hostname}

		sslSymlinkVhosts, err := repo.getSymlinkSslPairVhostsByVhost(vhost.Hostname)
		if err != nil {
			log.Printf("FailedToGetSymlinkSslPairVhosts (%s): %s", vhost.Hostname.String(), err.Error())
			continue
		}
		sslVhosts = append(sslVhosts, sslSymlinkVhosts...)

		vhostCertKeyContentStr, err := infraHelper.GetFileContent(vhostCertKeyFilePath)
		if err != nil {
			log.Printf("FailedToOpenCertKeyFile (%s): %s", hostnameStr, err.Error())
			continue
		}
		privateKey, err := valueObject.NewSslPrivateKey(vhostCertKeyContentStr)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), hostnameStr)
			continue
		}

		vhostCertFilePath := "/app/conf/pki/" + hostnameStr + ".crt"
		vhostCertFileContentStr, err := infraHelper.GetFileContent(vhostCertFilePath)
		if err != nil {
			log.Printf("FailedToOpenCertFile (%s): %s", hostnameStr, err.Error())
			continue
		}
		certificate, err := valueObject.NewSslCertificateContent(vhostCertFileContentStr)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), hostnameStr)
			continue
		}

		sslCertificates, err := repo.sslCertificatesFactory(certificate)
		if err != nil {
			log.Printf("FailedToGetMainAndChainedCerts (%s): %s", hostnameStr, err.Error())
			continue
		}

		ssl, err := repo.sslPairFactory(sslVhosts, privateKey, sslCertificates)
		if err != nil {
			log.Printf("FailedToGetSslPair (%s): %s", hostnameStr, err.Error())
			continue
		}

		sslPairs = append(sslPairs, ssl)
	}

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

func (repo SslQueryRepo) GetSslPairByVirtualHost(
	virtualHost valueObject.Fqdn,
) (entity.SslPair, error) {
	sslPairs, err := repo.GetSslPairs()
	if err != nil {
		return entity.SslPair{}, err
	}

	if len(sslPairs) < 1 {
		return entity.SslPair{}, errors.New("SslPairNotFound")
	}

	for _, sslPair := range sslPairs {
		for _, vhost := range sslPair.VirtualHosts {
			if vhost.String() != virtualHost.String() {
				continue
			}

			return sslPair, nil
		}
	}

	return entity.SslPair{}, errors.New("SslPairNotFound")
}
