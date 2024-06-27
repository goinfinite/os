package sslInfra

import (
	"errors"
	"log"
	"net"
	"os"

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

func (repo *SslCmdRepo) ReplaceWithSelfSigned(vhostName valueObject.Fqdn) error {
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(repo.persistentDbSvc)
	aliases, err := vhostQueryRepo.ReadAliasesByParentHostname(vhostName)
	if err != nil {
		return errors.New("ReadVhostAliasesError: " + err.Error())
	}

	aliasesHostname := []string{}
	for _, alias := range aliases {
		aliasesHostname = append(aliasesHostname, alias.Hostname.String())
	}

	err = repo.deleteCurrentSsl(vhostName)
	if err != nil {
		return err
	}

	return infraHelper.CreateSelfSignedSsl(
		infraData.GlobalConfigs.PkiConfDir,
		vhostName.String(),
		aliasesHostname,
	)
}

func (repo *SslCmdRepo) createDomainValidationMapping(
	targetVhostName valueObject.Fqdn,
	expectedOwnershipHash valueObject.Hash,
) (mappingId valueObject.MappingId, err error) {
	path, _ := valueObject.NewMappingPath(
		infraData.GlobalConfigs.DomainOwnershipValidationUrlPath,
	)
	matchPattern, _ := valueObject.NewMappingMatchPattern("equals")
	targetType, _ := valueObject.NewMappingTargetType("inline-html")
	httpResponseCode, _ := valueObject.NewHttpResponseCode(200)
	targetValue, _ := valueObject.NewMappingTargetValue(
		expectedOwnershipHash.String(), targetType,
	)

	inlineHtmlMapping := dto.NewCreateMapping(
		targetVhostName,
		path,
		matchPattern,
		targetType,
		&targetValue,
		&httpResponseCode,
	)

	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(repo.persistentDbSvc)
	mappingId, err = mappingCmdRepo.Create(inlineHtmlMapping)
	if err != nil {
		return mappingId, errors.New(
			"CreateOwnershipValidationMappingError (" +
				targetVhostName.String() + "): " + err.Error(),
		)
	}

	return mappingId, nil
}

func (repo *SslCmdRepo) filterDomainsMappedToSomewhere(
	vhostNames []valueObject.Fqdn,
) []valueObject.Fqdn {
	domainsMappedToSomewhere := []valueObject.Fqdn{}
	for _, vhostName := range vhostNames {
		vhostNameStr := vhostName.String()

		googleDnsServerIp := "8.8.8.8"
		rawVhostIps, err := infraHelper.RunCmdWithSubShell(
			"dig +short " + vhostNameStr + " @ " + googleDnsServerIp,
		)
		if err != nil || rawVhostIps == "" {
			cleanbrowsingSecurityFilterDnsServerIp := "185.228.168.9"
			rawVhostIps, err := infraHelper.RunCmdWithSubShell(
				"dig +short " + vhostNameStr + " @ " + cleanbrowsingSecurityFilterDnsServerIp,
			)
			if err != nil || rawVhostIps == "" {
				continue
			}
		}

		domainsMappedToSomewhere = append(domainsMappedToSomewhere, vhostName)
	}

	return domainsMappedToSomewhere
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

func (repo *SslCmdRepo) filterDomainsMappedToServer(
	vhostNames []valueObject.Fqdn,
	expectedOwnershipHash valueObject.Hash,
) []valueObject.Fqdn {
	domainsMappedToServer := []valueObject.Fqdn{}

	for _, vhostName := range vhostNames {
		vhostNameStr := vhostName.String()
		vhostNameIps, err := net.LookupIP(vhostNameStr)
		if err != nil || len(vhostNameIps) == 0 {
			continue
		}

		var serverIpAddress *valueObject.IpAddress
		for _, vhostNameIp := range vhostNameIps {
			ipAddress, err := valueObject.NewIpAddress(vhostNameIp.String())
			if err != nil {
				continue
			}

			serverIpAddress = &ipAddress
			break
		}

		if serverIpAddress == nil {
			continue
		}

		domainValidationMappingId, err := repo.createDomainValidationMapping(
			vhostName,
			expectedOwnershipHash,
		)
		if err != nil {
			continue
		}

		hashUrlPath := infraData.GlobalConfigs.DomainOwnershipValidationUrlPath
		hashUrlFull := "https://" + vhostNameStr + hashUrlPath
		serverIpStr := serverIpAddress.String()

		ownershipHashFound, err := infraHelper.RunCmd(
			"curl", "-skLm", "10", "--resolve", vhostNameStr+":443:"+serverIpStr, hashUrlFull,
		)
		if err != nil {
			hashUrlFull = "https://" + serverIpStr + hashUrlPath
			ownershipHashFound, err = infraHelper.RunCmd(
				"curl", "-skLm", "10", "-H", "Host: "+vhostNameStr, hashUrlFull,
			)
			if err != nil {
				continue
			}
		}

		if ownershipHashFound != expectedOwnershipHash.String() {
			continue
		}

		domainsMappedToServer = append(domainsMappedToServer, vhostName)

		if repo.shouldIncludeWww(vhostName) {
			vhostNameWithWww, err := valueObject.NewFqdn("www." + vhostNameStr)
			if err != nil {
				continue
			}

			domainsMappedToServer = append(domainsMappedToServer, vhostNameWithWww)
		}

		mappingCmdRepo := mappingInfra.NewMappingCmdRepo(repo.persistentDbSvc)
		err = mappingCmdRepo.Delete(domainValidationMappingId)
		if err != nil {
			log.Printf("DeleteOwnershipValidationMappingError: %s", err.Error())
		}
	}

	return domainsMappedToServer
}

func (repo *SslCmdRepo) ReplaceWithValidSsl(sslPair entity.SslPair) error {
	validPairVhostNames := repo.filterDomainsMappedToSomewhere(
		sslPair.VirtualHostsHostnames,
	)
	if len(validPairVhostNames) == 0 {
		return errors.New("NoDomainsMappedToSomewhere")
	}

	expectedOwnershipHash, err := repo.sslQueryRepo.GetOwnershipValidationHash(
		sslPair.Certificate.CertificateContent,
	)
	if err != nil {
		return errors.New(
			"CreateOwnershipValidationHashError: " + err.Error(),
		)
	}
	pairVhostNamesMappedToServer := repo.filterDomainsMappedToServer(
		validPairVhostNames, expectedOwnershipHash,
	)
	if len(pairVhostNamesMappedToServer) == 0 {
		return errors.New("NoDomainsMappedToServer")
	}

	firstVhostName := pairVhostNamesMappedToServer[0]
	firstVhostNameStr := firstVhostName.String()
	vhostRootDir := infraData.GlobalConfigs.PrimaryPublicDir
	if !infraHelper.IsPrimaryVirtualHost(firstVhostName) {
		vhostRootDir += "/" + firstVhostNameStr
	}

	certbotCmd := "certbot certonly --webroot --webroot-path " + vhostRootDir +
		" --agree-tos --register-unsafely-without-email --cert-name " + firstVhostNameStr

	for _, pairVhostName := range pairVhostNamesMappedToServer {
		certbotCmd += " -d " + pairVhostName.String()
	}

	_, err = infraHelper.RunCmdWithSubShell(certbotCmd)
	if err != nil {
		return errors.New("CreateValidSslError: " + err.Error())
	}

	certbotDirPath := "/etc/letsencrypt/live"
	shouldOverwrite := true

	certbotCrtFilePath := certbotDirPath + "/" + firstVhostNameStr + "/fullchain.pem"
	vhostCrtFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + firstVhostNameStr + ".crt"
	err = infraHelper.CreateSymlink(
		certbotCrtFilePath,
		vhostCrtFilePath,
		shouldOverwrite,
	)
	if err != nil {
		return errors.New("CreateSslCrtSymlinkError: " + err.Error())
	}

	certbotKeyFilePath := certbotDirPath + "/" + firstVhostNameStr + "/privkey.pem"
	vhostKeyFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + firstVhostNameStr + ".key"
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
	if len(createSslPair.VirtualHostsHostnames) == 0 {
		return errors.New("EmptyVirtualHosts")
	}

	firstVhostNameStr := createSslPair.VirtualHostsHostnames[0].String()
	firstVhostCertFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + firstVhostNameStr + ".crt"
	firstVhostCertKeyFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + firstVhostNameStr + ".key"

	for _, vhostName := range createSslPair.VirtualHostsHostnames {
		vhostStr := vhostName.String()
		vhostCertFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + vhostStr + ".crt"
		vhostCertKeyFilePath := infraData.GlobalConfigs.PkiConfDir + "/" + vhostStr + ".key"

		shouldBeSymlink := vhostStr != firstVhostNameStr
		if shouldBeSymlink {
			shouldOverwrite := true
			err := infraHelper.CreateSymlink(
				firstVhostCertFilePath,
				vhostCertFilePath,
				shouldOverwrite,
			)
			if err != nil {
				log.Printf(
					"CreateSslCertSymlinkError (%s): %s", vhostName.String(), err.Error(),
				)
				continue
			}

			err = infraHelper.CreateSymlink(
				firstVhostCertKeyFilePath,
				vhostCertKeyFilePath,
				shouldOverwrite,
			)
			if err != nil {
				log.Printf(
					"CreateSslKeySymlinkError (%s): %s", vhostName.String(), err.Error(),
				)
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

	for _, vhostName := range sslPairToDelete.VirtualHostsHostnames {
		err = repo.ReplaceWithSelfSigned(vhostName)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), vhostName.String())
			continue
		}
	}

	return nil
}

func (repo *SslCmdRepo) DeleteSslPairVhosts(
	deleteDto dto.DeleteSslPairVhosts,
) error {
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(repo.persistentDbSvc)
	for _, vhostName := range deleteDto.VirtualHostsHostnames {
		_, err := vhostQueryRepo.ReadByHostname(vhostName)
		if err != nil {
			continue
		}

		err = repo.ReplaceWithSelfSigned(vhostName)
		if err != nil {
			log.Printf(
				"DeleteSslPairVhostsError (%s): %s",
				vhostName.String(),
				err.Error(),
			)
		}
	}

	return nil
}
