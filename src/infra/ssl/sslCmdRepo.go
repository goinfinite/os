package sslInfra

import (
	"errors"
	"log"
	"net"
	"os"
	"reflect"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	envDataInfra "github.com/speedianet/os/src/infra/shared"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

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

	vhostCertFilePath := envDataInfra.PkiConfDir + "/" + vhostStr + ".crt"
	vhostCertFileExists := infraHelper.FileExists(vhostCertFilePath)
	if vhostCertFileExists {
		err := os.Remove(vhostCertFilePath)
		if err != nil {
			return errors.New("FailedToDeleteCertFile: " + err.Error())
		}
	}

	vhostCertKeyFilePath := envDataInfra.PkiConfDir + "/" + vhostStr + ".key"
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

	return infraHelper.CreateSelfSignedSsl(envDataInfra.PkiConfDir, vhost.String())
}

func (repo SslCmdRepo) shouldIncludeWww(host valueObject.Fqdn) bool {
	rootDomain, err := infraHelper.GetRootDomain(host)
	if err != nil {
		return false
	}

	hostStr := host.String()
	isSubdomain := rootDomain.String() != hostStr
	if isSubdomain {
		return false
	}

	wwwVhost := "www." + hostStr
	wwwVhostIps, err := net.LookupIP(wwwVhost)
	if err != nil {
		return false
	}

	wwwDnsExists := len(wwwVhostIps) > 0
	if !wwwDnsExists {
		return false
	}

	vhostIps, err := net.LookupIP(hostStr)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(wwwVhostIps, vhostIps)
}

func (repo SslCmdRepo) ReplaceWithValidSsl(sslPair entity.SslPair) error {
	path, _ := valueObject.NewMappingPath(envDataInfra.OwnershipValidationPath)
	matchPattern, _ := valueObject.NewMappingMatchPattern("equals")
	targetType, _ := valueObject.NewMappingTargetType("inline-html")
	httpResponseCode, _ := valueObject.NewHttpResponseCode(200)

	ownershipValidationHash := repo.sslQueryRepo.GetOwnershipValidationHash(
		sslPair.Certificate.CertificateContent,
	)
	inlineHtmlContent, _ := valueObject.NewInlineHtmlContent(ownershipValidationHash)

	firstVhost := sslPair.VirtualHosts[0]
	inlineHmtlMapping := dto.NewCreateMapping(
		firstVhost,
		path,
		matchPattern,
		targetType,
		nil,
		nil,
		&httpResponseCode,
		&inlineHtmlContent,
	)

	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}
	err := vhostCmdRepo.CreateMapping(inlineHmtlMapping)
	if err != nil {
		return errors.New("FailedToCreateOwnershipValidationMapping: " + err.Error())
	}

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	isOwnershipValid := vhostQueryRepo.CheckDomainOwnership(
		firstVhost,
		ownershipValidationHash,
	)

	vhostMappings, err := vhostQueryRepo.GetMappingsByHostname(firstVhost)
	if err != nil {
		return errors.New("FailedToGetVhostMappings: " + err.Error())
	}

	firstVhostStr := firstVhost.String()
	if len(vhostMappings) == 0 {
		return errors.New("VhostMappingsNotFound: " + firstVhostStr)
	}

	lastMappingIndex := len(vhostMappings) - 1
	lastMapping := vhostMappings[lastMappingIndex]

	err = vhostCmdRepo.DeleteMapping(lastMapping)
	if err != nil {
		return errors.New("FailedToDeleteOwnershipValidationMapping: " + err.Error())
	}

	if !isOwnershipValid {
		return errors.New("CurrentHostIsNotDomainOwner: " + firstVhostStr)
	}

	vhostRootDir := "/app/html"
	if !infraHelper.IsVirtualHostPrimaryDomain(firstVhost) {
		vhostRootDir += "/" + firstVhostStr
	}

	certbotCmd := "certbot certonly --webroot --webroot-path " + vhostRootDir +
		" --agree-tos --register-unsafely-without-email --cert-name " + firstVhostStr +
		" -d " + firstVhostStr

	shouldIncludeWww := repo.shouldIncludeWww(firstVhost)
	if shouldIncludeWww {
		certbotCmd += " -d www." + firstVhostStr
	}

	_, err = infraHelper.RunCmdWithSubShell(certbotCmd)
	if err != nil {
		return errors.New("CreateValidSslFailed: " + err.Error())
	}

	certbotDirPath := "/etc/letsencrypt/live"
	shouldOverwrite := true

	certbotCrtFilePath := certbotDirPath + "/" + firstVhostStr + "/fullchain.pem"
	vhostCrtFilePath := envDataInfra.PkiConfDir + "/" + firstVhostStr + ".crt"
	err = infraHelper.CreateSymlink(
		certbotCrtFilePath,
		vhostCrtFilePath,
		shouldOverwrite,
	)
	if err != nil {
		return errors.New("CreateSslCrtSymlinkError: " + err.Error())
	}

	certbotKeyFilePath := certbotDirPath + "/" + firstVhostStr + "/privkey.pem"
	vhostKeyFilePath := envDataInfra.PkiConfDir + "/" + firstVhostStr + ".key"
	err = infraHelper.CreateSymlink(
		certbotKeyFilePath,
		vhostKeyFilePath,
		shouldOverwrite,
	)
	if err != nil {
		return errors.New("CreateSslKeySymlinkError: " + err.Error())
	}

	return nil
}

func (repo SslCmdRepo) Create(createSslPair dto.CreateSslPair) error {
	if len(createSslPair.VirtualHosts) == 0 {
		return errors.New("NoVirtualHostsProvidedToCreateSslPair")
	}

	firstVhostStr := createSslPair.VirtualHosts[0].String()
	firstVhostCertFilePath := envDataInfra.PkiConfDir + "/" + firstVhostStr + ".crt"
	firstVhostCertKeyFilePath := envDataInfra.PkiConfDir + "/" + firstVhostStr + ".key"

	for _, vhost := range createSslPair.VirtualHosts {
		vhostStr := vhost.String()
		vhostCertFilePath := envDataInfra.PkiConfDir + "/" + vhostStr + ".crt"
		vhostCertKeyFilePath := envDataInfra.PkiConfDir + "/" + vhostStr + ".key"

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
