package sslInfra

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	infraData "github.com/speedianet/os/src/infra/infraData"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
)

type SslCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	sslQueryRepo    SslQueryRepo
}

func NewSslCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *SslCmdRepo {
	return &SslCmdRepo{
		persistentDbSvc: persistentDbSvc,
		sslQueryRepo:    SslQueryRepo{},
	}
}

func (repo *SslCmdRepo) deleteCurrentSsl(vhost valueObject.Fqdn) error {
	vhostStr := vhost.String()

	vhostCertFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + vhostStr + ".crt"
	vhostCertFileExists := infraHelper.FileExists(vhostCertFilePath)
	if vhostCertFileExists {
		err := os.Remove(vhostCertFilePath)
		if err != nil {
			return errors.New("DeleteCertFileError: " + err.Error())
		}
	}

	vhostCertKeyFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + vhostStr + ".key"
	vhostCertKeyFileExists := infraHelper.FileExists(vhostCertKeyFilePath)
	if vhostCertKeyFileExists {
		err := os.Remove(vhostCertKeyFilePath)
		if err != nil {
			return errors.New("DeleteCertKeyFileError: " + err.Error())
		}
	}

	return nil
}

func (repo *SslCmdRepo) ReplaceWithSelfSigned(vhost valueObject.Fqdn) error {
	err := repo.deleteCurrentSsl(vhost)
	if err != nil {
		return err
	}

	return infraHelper.CreateSelfSignedSsl(infraData.GlobalConfigs.PkiConfDir, vhost.String())
}

func (repo *SslCmdRepo) isDomainMappedToServer(
	vhost valueObject.Fqdn,
	expectedOwnershipHash valueObject.Hash,
) bool {
	vhostStr := vhost.String()

	rawVhostIps, err := infraHelper.RunCmd("dig", "+short", vhostStr, "@8.8.8.8")
	if err != nil || rawVhostIps == "" {
		rawVhostIps, err = infraHelper.RunCmd("dig", "+short", vhostStr, "@1.1.1.1")
		if err != nil || rawVhostIps == "" {
			return false
		}
	}

	rawVhostIpsParts := strings.Split(rawVhostIps, "\n")
	if len(rawVhostIpsParts) == 0 {
		return false
	}

	var serverIpAddress *valueObject.IpAddress
	for _, rawVhostIp := range rawVhostIpsParts {
		ipAddress, err := valueObject.NewIpAddress(rawVhostIp)
		if err != nil {
			continue
		}

		serverIpAddress = &ipAddress
		break
	}

	if serverIpAddress == nil {
		return false
	}

	ownershipValidateUrl := "https://" + serverIpAddress.String() +
		infraData.GlobalConfigs.DomainOwnershipValidationUrlPath

	ownershipHashFound, err := infraHelper.RunCmd(
		"curl",
		"-skL",
		"--max-time",
		"10",
		"--header",
		"Host: "+vhostStr,
		ownershipValidateUrl,
	)
	if err != nil {
		return false
	}

	return ownershipHashFound == expectedOwnershipHash.String()
}

func (repo *SslCmdRepo) shouldIncludeWww(vhost valueObject.Fqdn) bool {
	rootDomain, err := infraHelper.GetRootDomain(vhost)
	if err != nil {
		return false
	}

	vhostStr := vhost.String()
	isSubdomain := rootDomain.String() != vhostStr
	if isSubdomain {
		return false
	}

	wwwDnsEntry := "www." + vhostStr
	wwwDnsEntryIps, err := net.LookupIP(wwwDnsEntry)
	if err != nil {
		return false
	}

	wwwDnsEntryExists := len(wwwDnsEntryIps) > 0
	if !wwwDnsEntryExists {
		return false
	}

	vhostIps, err := net.LookupIP(vhostStr)
	if err != nil {
		return false
	}

	firstVhostIp := vhostIps[0]
	for _, wwwDnsEntryIp := range wwwDnsEntryIps {
		if !firstVhostIp.Equal(wwwDnsEntryIp) {
			continue
		}

		return true
	}

	return false
}

func (repo *SslCmdRepo) ReplaceWithValidSsl(sslPair entity.SslPair) error {
	path, _ := valueObject.NewMappingPath(infraData.GlobalConfigs.DomainOwnershipValidationUrlPath)
	matchPattern, _ := valueObject.NewMappingMatchPattern("equals")
	targetType, _ := valueObject.NewMappingTargetType("inline-html")
	httpResponseCode, _ := valueObject.NewHttpResponseCode(200)

	expectedOwnershipHash, err := repo.sslQueryRepo.GetOwnershipValidationHash(
		sslPair.Certificate.CertificateContent,
	)
	if err != nil {
		return errors.New("CreateOwnershipValidationHashError: " + err.Error())
	}
	targetValue, _ := valueObject.NewMappingTargetValue(
		expectedOwnershipHash.String(), targetType,
	)

	firstVhost := sslPair.VirtualHosts[0]
	inlineHtmlMapping := dto.NewCreateMapping(
		firstVhost,
		path,
		matchPattern,
		targetType,
		&targetValue,
		&httpResponseCode,
	)

	firstVhostStr := firstVhost.String()
	ips, err := net.LookupHost(firstVhostStr)
	if err != nil {
		return errors.New("LookupHostError: " + err.Error())
	}

	if len(ips) == 0 {
		return errors.New("VhostDoesNotPointToAnyIp")
	}

	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(repo.persistentDbSvc)
	mappingId, err := mappingCmdRepo.Create(inlineHtmlMapping)
	if err != nil {
		return errors.New("CreateOwnershipValidationMappingError: " + err.Error())
	}

	isDomainMappedToServer := repo.isDomainMappedToServer(
		firstVhost,
		expectedOwnershipHash,
	)

	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(repo.persistentDbSvc)
	mappings, err := mappingQueryRepo.ReadByHostname(firstVhost)
	if err != nil {
		return errors.New("ReadVhostMappingsError: " + err.Error())
	}

	if len(mappings) == 0 {
		return errors.New("VhostMappingsNotFound")
	}

	err = mappingCmdRepo.Delete(mappingId)
	if err != nil {
		return errors.New("DeleteOwnershipValidationMappingError: " + err.Error())
	}

	if !isDomainMappedToServer {
		return errors.New("DomainNotResolvingToServer")
	}

	vhostRootDir := infraData.GlobalConfigs.PrimaryPublicDir
	if !infraHelper.IsPrimaryVirtualHost(firstVhost) {
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
		return errors.New("CreateValidSslError: " + err.Error())
	}

	certbotDirPath := "/etc/letsencrypt/live"
	shouldOverwrite := true

	certbotCrtFilePath := certbotDirPath + "/" + firstVhostStr + "/fullchain.pem"
	vhostCrtFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + firstVhostStr + ".crt"
	err = infraHelper.CreateSymlink(
		certbotCrtFilePath,
		vhostCrtFilePath,
		shouldOverwrite,
	)
	if err != nil {
		return errors.New("CreateSslCrtSymlinkError: " + err.Error())
	}

	certbotKeyFilePath := certbotDirPath + "/" + firstVhostStr + "/privkey.pem"
	vhostKeyFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + firstVhostStr + ".key"
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

func (repo *SslCmdRepo) Create(createSslPair dto.CreateSslPair) error {
	if len(createSslPair.VirtualHosts) == 0 {
		return errors.New("EmptyVirtualHosts")
	}

	firstVhostStr := createSslPair.VirtualHosts[0].String()
	firstVhostCertFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + firstVhostStr + ".crt"
	firstVhostCertKeyFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + firstVhostStr + ".key"

	for _, vhost := range createSslPair.VirtualHosts {
		vhostStr := vhost.String()
		vhostCertFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + vhostStr + ".crt"
		vhostCertKeyFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + vhostStr + ".key"

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

func (repo *SslCmdRepo) Delete(sslId valueObject.SslId) error {
	sslPairToDelete, err := repo.sslQueryRepo.ReadById(sslId)
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

func (repo *SslCmdRepo) DeleteSslPairVhosts(
	deleteDto dto.DeleteSslPairVhosts,
) error {
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	for _, vhost := range deleteDto.VirtualHosts {
		_, err := vhostQueryRepo.GetByHostname(vhost)
		if err != nil {
			continue
		}

		err = repo.ReplaceWithSelfSigned(vhost)
		if err != nil {
			log.Printf(
				"DeleteSslPairVhostsError (%s): %s",
				vhost.String(),
				err.Error(),
			)
		}
	}

	return nil
}
