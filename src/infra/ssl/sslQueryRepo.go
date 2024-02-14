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
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

const nginxConfDir = "/app/conf/nginx"

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

func (repo SslQueryRepo) getHostnameFromCertFilePath(
	certFilePath string,
) (valueObject.Fqdn, error) {
	certFilePathWithoutBase := filepath.Base(certFilePath)
	certFileNameWithoutExt := strings.TrimSuffix(certFilePathWithoutBase, ".crt")

	return valueObject.NewFqdn(certFileNameWithoutExt)
}

func (repo SslQueryRepo) getCertFilesPathWithSymlinkVhosts(
	pkiFiles []entity.UnixFile,
) (map[string][]valueObject.Fqdn, error) {
	certFilesPathWithSymlinkVhostsMap := map[string][]valueObject.Fqdn{}
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

		_, targetCertFilePathAlreadyExistsInMap := certFilesPathWithSymlinkVhostsMap[certFilePathStr]
		if !targetCertFilePathAlreadyExistsInMap {
			certFilesPathWithSymlinkVhostsMap[targetCertFilePathStr] = []valueObject.Fqdn{}
		}

		if certFilePathStr == targetCertFilePathStr {
			continue
		}

		symlinkVhost, err := repo.getHostnameFromCertFilePath(certFilePathStr)
		if err != nil {
			log.Printf("%s: %s", err.Error(), certFilePathStr)
			continue
		}

		certFilesPathWithSymlinkVhostsMap[targetCertFilePathStr] = append(
			certFilesPathWithSymlinkVhostsMap[targetCertFilePathStr],
			symlinkVhost,
		)
	}

	return certFilesPathWithSymlinkVhostsMap, nil
}

func (repo SslQueryRepo) GetSslPairs() ([]entity.SslPair, error) {
	sslPairs := []entity.SslPair{}

	pkiFilesPath, _ := valueObject.NewUnixFilePath("/app/conf/pki")

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	pkiFiles, err := filesQueryRepo.Get(pkiFilesPath)
	if err != nil {
		return sslPairs, errors.New("FailedToGetPkiFiles: " + err.Error())
	}

	certFilesPathWithSymlinkVhostsMap, err := repo.getCertFilesPathWithSymlinkVhosts(pkiFiles)
	if err != nil {
		return sslPairs, errors.New("FailedToGetCertFilesPathWithSymlinksVhosts: " + err.Error())
	}

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	for targetCertFilePath, symlinkVhost := range certFilesPathWithSymlinkVhostsMap {
		targetVhost, err := repo.getHostnameFromCertFilePath(targetCertFilePath)
		if err != nil {
			log.Printf("%s: %s", err.Error(), targetCertFilePath)
			continue
		}

		targetVhostStr := targetVhost.String()

		vhostConfigFilePath, err := vhostQueryRepo.GetVirtualHostConfFilePath(targetVhost)
		if err != nil {
			log.Printf("FailedToGetVhostConfFile (%s): %s", targetVhostStr, err.Error())
		}

		vhostConfigContentStr, err := infraHelper.GetFileContent(vhostConfigFilePath.String())
		if err != nil {
			log.Printf("FailedToOpenVhostConfFile (%s): %s", targetVhostStr, err.Error())
			continue
		}

		fileIsEmpty := len(vhostConfigContentStr) < 1
		if fileIsEmpty {
			log.Printf("VirtualHostConfFileIsEmpty (%s)", targetVhostStr)
			continue
		}

		vhostCertKeyFilePath := "/app/conf/pki/" + targetVhostStr + ".key"
		vhostCertKeyContentStr, err := infraHelper.GetFileContent(vhostCertKeyFilePath)
		if err != nil {
			log.Printf("FailedToOpenCertKeyFile (%s): %s", targetVhostStr, err.Error())
			continue
		}
		privateKey, err := valueObject.NewSslPrivateKey(vhostCertKeyContentStr)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), targetVhostStr)
			continue
		}

		vhostCertFilePath := "/app/conf/pki/" + targetVhostStr + ".crt"
		vhostCertFileContentStr, err := infraHelper.GetFileContent(vhostCertFilePath)
		if err != nil {
			log.Printf("FailedToOpenCertFile (%s): %s", targetVhostStr, err.Error())
			continue
		}
		certificate, err := valueObject.NewSslCertificateContent(vhostCertFileContentStr)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), targetVhostStr)
			continue
		}

		sslCertificates, err := repo.sslCertificatesFactory(certificate)
		if err != nil {
			log.Printf("FailedToGetMainAndChainedCerts (%s): %s", targetVhostStr, err.Error())
			continue
		}

		sslVhosts := symlinkVhost
		sslVhosts = append(sslVhosts, targetVhost)

		ssl, err := repo.sslPairFactory(sslVhosts, privateKey, sslCertificates)
		if err != nil {
			log.Printf("FailedToGetSslPair (%s): %s", targetVhostStr, err.Error())
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
