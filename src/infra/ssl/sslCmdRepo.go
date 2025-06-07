package sslInfra

import (
	"errors"
	"log/slog"
	"os"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	o11yInfra "github.com/goinfinite/os/src/infra/o11y"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

const DomainOwnershipValidationUrlPath string = "/validateOwnership"

type SslCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	sslQueryRepo    *SslQueryRepo
	vhostQueryRepo  *vhostInfra.VirtualHostQueryRepo
}

func NewSslCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *SslCmdRepo {
	return &SslCmdRepo{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		sslQueryRepo:    &SslQueryRepo{},
		vhostQueryRepo:  vhostInfra.NewVirtualHostQueryRepo(persistentDbSvc),
	}
}

func (repo *SslCmdRepo) dnsFilterFunctionalHostnames(
	vhostHostnames []valueObject.Fqdn,
	serverPublicIpAddress valueObject.IpAddress,
) []valueObject.Fqdn {
	functionalHostnames := []valueObject.Fqdn{}

	for _, vhostHostname := range vhostHostnames {
		wwwVirtualHostHostname, err := valueObject.NewFqdn("www." + vhostHostname.String())
		if err != nil {
			slog.Debug(
				"InvalidWwwVirtualHostHostname",
				slog.String("fqdn", vhostHostname.String()),
			)
			continue
		}

		vhostHostnames = append(vhostHostnames, wwwVirtualHostHostname)
	}

	serverPublicIpAddressStr := serverPublicIpAddress.String()
	for _, vhostHostname := range vhostHostnames {
		vhostHostnameStr := vhostHostname.String()

		hostnameRecords, err := infraHelper.DnsLookup(vhostHostnameStr, nil)
		if err != nil {
			slog.Debug(
				"DnsLookupFailed",
				slog.String("fqdn", vhostHostnameStr),
				slog.String("error", err.Error()),
			)
			continue
		}

		if len(hostnameRecords) == 0 {
			slog.Debug("NoDnsRecordsFound", slog.String("fqdn", vhostHostnameStr))
			continue
		}

		for _, dnsRecord := range hostnameRecords {
			if dnsRecord != serverPublicIpAddressStr {
				slog.Debug(
					"DnsRecordDoesNotMatchServerIpAddress",
					slog.String("fqdn", vhostHostnameStr),
					slog.String("dnsRecord", dnsRecord),
					slog.String("serverPublicIpAddress", serverPublicIpAddressStr),
				)
				continue
			}

			functionalHostnames = append(functionalHostnames, vhostHostname)
			break
		}
	}

	return functionalHostnames
}

func (repo *SslCmdRepo) createOwnershipValidationMapping(
	mappingCmdRepo *vhostInfra.MappingCmdRepo,
	targetVirtualHostHostname valueObject.Fqdn,
	expectedOwnershipHash valueObject.Hash,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) (mappingId valueObject.MappingId, err error) {
	path, _ := valueObject.NewMappingPath(DomainOwnershipValidationUrlPath)
	matchPattern, _ := valueObject.NewMappingMatchPattern("equals")
	targetType, _ := valueObject.NewMappingTargetType("inline-html")
	httpResponseCode, _ := valueObject.NewHttpResponseCode(200)
	targetValue, _ := valueObject.NewMappingTargetValue(
		expectedOwnershipHash.String(), targetType,
	)
	shouldUpgradeInsecureRequests := false

	inlineHtmlMapping := dto.NewCreateMapping(
		targetVirtualHostHostname, path, matchPattern, targetType, &targetValue,
		&httpResponseCode, &shouldUpgradeInsecureRequests, nil,
		operatorAccountId, operatorIpAddress,
	)

	return mappingCmdRepo.Create(inlineHtmlMapping)
}

func (repo *SslCmdRepo) httpFilterFunctionalHostnames(
	vhostHostnames []valueObject.Fqdn,
	expectedOwnershipHash valueObject.Hash,
	serverPublicIpAddress valueObject.IpAddress,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) []valueObject.Fqdn {
	functionalHostnames := []valueObject.Fqdn{}

	serverPublicIpAddressStr := serverPublicIpAddress.String()
	expectedHashStr := expectedOwnershipHash.String()
	mappingCmdRepo := vhostInfra.NewMappingCmdRepo(repo.persistentDbSvc)

	for _, vhostHostname := range vhostHostnames {
		vhostHostnameStr := vhostHostname.String()
		ownershipValidationMappingId, err := repo.createOwnershipValidationMapping(
			mappingCmdRepo, vhostHostname, expectedOwnershipHash, operatorAccountId,
			operatorIpAddress,
		)
		if err != nil {
			continue
		}

		hashUrlPath := DomainOwnershipValidationUrlPath
		hashUrlFull := "https://" + vhostHostnameStr + hashUrlPath
		curlBaseCmd := "curl -skLm 10 "
		sniFlag := "--resolve " + vhostHostnameStr + ":443:" + serverPublicIpAddressStr
		ownershipHashFound, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
			Command:               curlBaseCmd + sniFlag + " " + hashUrlFull,
			ShouldRunWithSubShell: true,
		})
		if err != nil {
			hashUrlFull = "https://" + serverPublicIpAddressStr + hashUrlPath
			ownershipHashFound, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
				Command:               curlBaseCmd + "-H \"Host: " + vhostHostnameStr + "\" " + hashUrlFull,
				ShouldRunWithSubShell: true,
			})
			if err != nil {
				continue
			}
		}

		if ownershipHashFound != expectedHashStr {
			continue
		}

		functionalHostnames = append(functionalHostnames, vhostHostname)

		err = mappingCmdRepo.Delete(ownershipValidationMappingId)
		if err != nil {
			slog.Error("DeleteOwnershipValidationMappingError", slog.String("error", err.Error()))
		}
	}

	return functionalHostnames
}

func (repo *SslCmdRepo) issueValidSsl(
	mainHostname valueObject.Fqdn,
	functionalHostnames []valueObject.Fqdn,
) error {
	mainHostnameStr := mainHostname.String()
	vhostRootDir := infraEnvs.PrimaryPublicDir
	if !infraHelper.IsPrimaryVirtualHost(mainHostname) {
		vhostRootDir += "/" + mainHostnameStr
	}

	if !infraHelper.FileExists(vhostRootDir) {
		return errors.New("VirtualHostRootDirNotFound")
	}

	certbotCmd := "certbot certonly --webroot --webroot-path " + vhostRootDir +
		" --agree-tos --register-unsafely-without-email --cert-name " + mainHostnameStr
	for _, functionalHostname := range functionalHostnames {
		certbotCmd += " -d " + functionalHostname.String()
	}

	_, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               certbotCmd,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return errors.New("GenerateValidSslCertError: " + err.Error())
	}

	certbotDirPath := "/etc/letsencrypt/live"
	shouldOverwrite := true

	certbotCrtFilePath := certbotDirPath + "/" + mainHostnameStr + "/fullchain.pem"
	vhostCrtFilePath := infraEnvs.PkiConfDir + "/" + mainHostnameStr + ".crt"
	err = infraHelper.CreateSymlink(certbotCrtFilePath, vhostCrtFilePath, shouldOverwrite)
	if err != nil {
		return errors.New("CreateSslCertSymlinkError: " + err.Error())
	}

	certbotKeyFilePath := certbotDirPath + "/" + mainHostnameStr + "/privkey.pem"
	vhostKeyFilePath := infraEnvs.PkiConfDir + "/" + mainHostnameStr + ".key"
	err = infraHelper.CreateSymlink(certbotKeyFilePath, vhostKeyFilePath, shouldOverwrite)
	if err != nil {
		return errors.New("CreateSslKeySymlinkError: " + err.Error())
	}

	return infraHelper.ReloadWebServer()
}

func (repo *SslCmdRepo) CreatePubliclyTrusted(
	createDto dto.CreatePubliclyTrustedSslPair,
) (sslPairId valueObject.SslPairId, err error) {
	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(repo.transientDbSvc)
	serverPublicIpAddress, err := o11yQueryRepo.ReadServerPublicIpAddress()
	if err != nil {
		return sslPairId, err
	}

	vhostReadResponse, err := repo.vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
		Pagination: dto.PaginationSingleItem,
		Hostname:   &createDto.VirtualHostHostname,
	})
	if err != nil {
		return sslPairId, errors.New("ReadVirtualHostEntitiesError: " + err.Error())
	}

	if len(vhostReadResponse.VirtualHosts) == 0 {
		return sslPairId, errors.New("VirtualHostNotFound")
	}

	virtualHostsHostnames := []valueObject.Fqdn{createDto.VirtualHostHostname}
	virtualHostsHostnames = append(
		virtualHostsHostnames, vhostReadResponse.VirtualHosts[0].AliasesHostnames...,
	)

	dnsFunctionalHostnames := repo.dnsFilterFunctionalHostnames(
		virtualHostsHostnames, serverPublicIpAddress,
	)
	if len(dnsFunctionalHostnames) == 0 {
		return sslPairId, errors.New("NoSslHostnamePointingToServerIpAddress")
	}

	synthesizer := tkInfra.Synthesizer{}
	dummyValue := synthesizer.PasswordFactory(32, false)
	dummyHash := infraHelper.GenStrongHash(dummyValue)

	expectedOwnershipHash, err := valueObject.NewHash(dummyHash)
	if err != nil {
		return sslPairId, errors.New("CreateOwnershipValidationHashError: " + err.Error())
	}

	httpFunctionalHostnames := repo.httpFilterFunctionalHostnames(
		dnsFunctionalHostnames, expectedOwnershipHash, serverPublicIpAddress,
		createDto.OperatorAccountId, createDto.OperatorIpAddress,
	)
	if len(httpFunctionalHostnames) == 0 {
		return sslPairId, errors.New("NoSslHostnamePassingHttpOwnershipValidation")
	}

	err = repo.issueValidSsl(createDto.VirtualHostHostname, httpFunctionalHostnames)
	if err != nil {
		return sslPairId, errors.New("IssueValidSslError: " + err.Error())
	}

	sslPairEntity, err := repo.sslQueryRepo.ReadFirst(dto.ReadSslPairsRequest{
		VirtualHostHostname: &createDto.VirtualHostHostname,
	})
	if err != nil {
		return sslPairId, errors.New("SslPairNotFound: " + err.Error())
	}

	return sslPairEntity.Id, nil
}

func (repo *SslCmdRepo) Create(
	createDto dto.CreateSslPair,
) (sslPairId valueObject.SslPairId, err error) {
	if len(createDto.VirtualHostsHostnames) == 0 {
		return sslPairId, errors.New("EmptyVirtualHosts")
	}

	for _, vhostHostname := range createDto.VirtualHostsHostnames {
		vhostHostnameStr := vhostHostname.String()
		vhostCertFilePath := infraEnvs.PkiConfDir + "/" + vhostHostnameStr + ".crt"
		vhostCertKeyFilePath := infraEnvs.PkiConfDir + "/" + vhostHostnameStr + ".key"

		certContentStr := createDto.Certificate.CertificateContent.String()
		if createDto.ChainCertificates != nil {
			certContentStr += "\n" + createDto.ChainCertificates.CertificateContent.String()
		}

		shouldOverwrite := true
		err := infraHelper.UpdateFile(vhostCertFilePath, certContentStr, shouldOverwrite)
		if err != nil {
			return sslPairId, errors.New("UpdateCertFileError: " + err.Error())
		}

		err = infraHelper.UpdateFile(
			vhostCertKeyFilePath, createDto.Key.String(), shouldOverwrite,
		)
		if err != nil {
			return sslPairId, errors.New("UpdateCertKeyFileError: " + err.Error())
		}
	}

	sslPairEntity, err := repo.sslQueryRepo.ReadFirst(dto.ReadSslPairsRequest{
		VirtualHostHostname: &createDto.VirtualHostsHostnames[0],
	})
	if err != nil {
		return sslPairId, errors.New("SslPairNotFound: " + err.Error())
	}

	err = infraHelper.ReloadWebServer()
	if err != nil {
		return sslPairId, errors.New("ReloadWebServerError: " + err.Error())
	}

	return sslPairEntity.Id, nil
}

func (repo *SslCmdRepo) ReplaceWithSelfSigned(vhostHostname valueObject.Fqdn) error {
	vhostEntity, err := repo.vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &vhostHostname,
	})
	if err != nil {
		return errors.New("ReadVirtualHostEntityError: " + err.Error())
	}
	if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
		return errors.New("AliasVirtualHostSslReliesOnParent")
	}

	aliasesVirtualHostsReadResponse, err := repo.vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
		Pagination:     dto.PaginationUnpaginated,
		ParentHostname: &vhostHostname,
	})
	if err != nil {
		return errors.New("ReadAliasesError: " + err.Error())
	}

	aliasesHostnames := []valueObject.Fqdn{}
	for _, aliasVirtualHostEntity := range aliasesVirtualHostsReadResponse.VirtualHosts {
		aliasesHostnames = append(aliasesHostnames, aliasVirtualHostEntity.Hostname)
	}

	vhostHostnameStr := vhostHostname.String()
	vhostCertFilePath := infraEnvs.PkiConfDir + "/" + vhostHostnameStr + ".crt"
	vhostCertFileExists := infraHelper.FileExists(vhostCertFilePath)
	if vhostCertFileExists {
		err := os.Remove(vhostCertFilePath)
		if err != nil {
			return errors.New("DeleteCertFileError: " + err.Error())
		}
	}

	vhostCertKeyFilePath := infraEnvs.PkiConfDir + "/" + vhostHostnameStr + ".key"
	vhostCertKeyFileExists := infraHelper.FileExists(vhostCertKeyFilePath)
	if vhostCertKeyFileExists {
		err := os.Remove(vhostCertKeyFilePath)
		if err != nil {
			return errors.New("DeleteCertKeyFileError: " + err.Error())
		}
	}

	pkiConfDir, err := valueObject.NewUnixFilePath(infraEnvs.PkiConfDir)
	if err != nil {
		return errors.New("PkiConfDirError: " + err.Error())
	}

	err = infraHelper.CreateSelfSignedSsl(pkiConfDir, vhostHostname, aliasesHostnames)
	if err != nil {
		return errors.New("CreateSelfSignedSslError: " + err.Error())
	}

	return infraHelper.ReloadWebServer()
}

func (repo *SslCmdRepo) Delete(sslPairId valueObject.SslPairId) error {
	sslPairEntity, err := repo.sslQueryRepo.ReadFirst(dto.ReadSslPairsRequest{
		SslPairId: &sslPairId,
	})
	if err != nil {
		return errors.New("ReadSslPairEntityError: " + err.Error())
	}

	err = repo.ReplaceWithSelfSigned(sslPairEntity.VirtualHostHostname)
	if err != nil {
		return errors.New("ReplaceWithSelfSignedError: " + err.Error())
	}

	return nil
}
