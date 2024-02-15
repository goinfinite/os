package sslInfra

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	filesInfra "github.com/speedianet/os/src/infra/files"
	infraHelper "github.com/speedianet/os/src/infra/helper"
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

	vhostCertKeyFilePath := pkiConfDir + firstVhostStr + ".key"
	vhostCertKeyContentStr, err := infraHelper.GetFileContent(vhostCertKeyFilePath)
	if err != nil {
		return ssl, errors.New("FailedToOpenCertKeyFile: " + err.Error())
	}
	privateKey, err := valueObject.NewSslPrivateKey(vhostCertKeyContentStr)
	if err != nil {
		return ssl, errors.New(err.Error() + "(" + firstVhostStr + ")")
	}

	vhostCertFilePath := pkiConfDir + firstVhostStr + ".crt"
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

func (repo SslQueryRepo) getHostnameFromCertFilePath(
	certFilePath string,
) (valueObject.Fqdn, error) {
	certFilePathWithoutBase := filepath.Base(certFilePath)
	certFileNameWithoutExt := strings.TrimSuffix(certFilePathWithoutBase, ".crt")

	return valueObject.NewFqdn(certFileNameWithoutExt)
}

func (repo SslQueryRepo) getCertFilesPathWithVhosts(
	pkiFiles []entity.UnixFile,
) (map[string][]valueObject.Fqdn, error) {
	certFilesPathWithVhostsMap := map[string][]valueObject.Fqdn{}
	for _, pkiFile := range pkiFiles {
		// TODO: remove when PR 19 is merged
		isDirPath := pkiFile.MimeType.IsDir()
		if isDirPath {
			continue
		}

		isCertFile := pkiFile.Extension.String() == "crt"
		if !isCertFile {
			continue
		}

		certFilePathStr := pkiFile.Path.String()
		targetCertFilePathStr := certFilePathStr

		isSymlink := infraHelper.IsSymlink(certFilePathStr)
		if isSymlink {
			targetCertFilePathFromSymlink, err := os.Readlink(certFilePathStr)
			if err != nil {
				log.Printf("FailedToGetCrtFilePathFromSymlink: %s", err.Error())
				continue
			}

			targetCertFilePathStr = targetCertFilePathFromSymlink
		}

		_, targetVhostAlreadyExistsInMap := certFilesPathWithVhostsMap[targetCertFilePathStr]
		if !targetVhostAlreadyExistsInMap {
			certFilesPathWithVhostsMap[targetCertFilePathStr] = []valueObject.Fqdn{}
		}

		vhost, err := repo.getHostnameFromCertFilePath(certFilePathStr)
		if err != nil {
			log.Printf("%s: %s", err.Error(), certFilePathStr)
			continue
		}

		certFilesPathWithVhostsMap[targetCertFilePathStr] = append(
			certFilesPathWithVhostsMap[targetCertFilePathStr],
			vhost,
		)
	}

	return certFilesPathWithVhostsMap, nil
}

func (repo SslQueryRepo) GetSslPairs() ([]entity.SslPair, error) {
	sslPairs := []entity.SslPair{}

	pkiFilesPath, _ := valueObject.NewUnixFilePath(pkiConfDir)

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	pkiFiles, err := filesQueryRepo.Get(pkiFilesPath)
	if err != nil {
		return sslPairs, errors.New("FailedToGetPkiFiles: " + err.Error())
	}

	certFilesPathWithVhostsMap, err := repo.getCertFilesPathWithVhosts(pkiFiles)
	if err != nil {
		return sslPairs, errors.New("FailedToGetCertFilesPathWithSymlinksVhosts: " + err.Error())
	}

	for certFilePath, sslVhosts := range certFilesPathWithVhostsMap {
		targetVhost, err := repo.getHostnameFromCertFilePath(certFilePath)
		if err != nil {
			log.Printf("%s: %s", err.Error(), certFilePath)
			continue
		}

		ssl, err := repo.sslPairFactory(sslVhosts)
		if err != nil {
			log.Printf("FailedToGetSslPair (%s): %s", targetVhost.String(), err.Error())
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
