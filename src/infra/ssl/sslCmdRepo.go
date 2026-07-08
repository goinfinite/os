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
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

const DomainOwnershipValidationUrlPath string = "/validateOwnership"

type SslCmdRepo struct {
	persistentDbSvc         *internalDbInfra.PersistentDatabaseService
	transientDbSvc          *internalDbInfra.TransientDatabaseService
	sslQueryRepo            *SslQueryRepo
	vhostHelpers            *vhostInfra.VirtualHostHelpers
	vhostQueryRepo          *vhostInfra.VirtualHostQueryRepo
	mappingCmdRepo          *vhostInfra.MappingCmdRepo
	mappingQueryRepo        *vhostInfra.MappingQueryRepo
	fileClerk               tkInfra.FileClerk
	ownershipValidationPath valueObject.MappingPath
}

func NewSslCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *SslCmdRepo {
	ownershipValidationPath, _ := valueObject.NewMappingPath(DomainOwnershipValidationUrlPath)
	return &SslCmdRepo{
		persistentDbSvc:         persistentDbSvc,
		transientDbSvc:          transientDbSvc,
		sslQueryRepo:            NewSslQueryRepo(),
		vhostHelpers:            vhostInfra.NewVirtualHostHelpers(),
		vhostQueryRepo:          vhostInfra.NewVirtualHostQueryRepo(persistentDbSvc),
		mappingCmdRepo:          vhostInfra.NewMappingCmdRepo(persistentDbSvc),
		mappingQueryRepo:        vhostInfra.NewMappingQueryRepo(persistentDbSvc),
		fileClerk:               tkInfra.FileClerk{},
		ownershipValidationPath: ownershipValidationPath,
	}
}

func (repo *SslCmdRepo) dnsFunctionalHostnamesFilter(
	vhostHostnames []tkValueObject.Fqdn,
	serverPublicIpAddress tkValueObject.IpAddress,
) []tkValueObject.Fqdn {
	functionalHostnames := []tkValueObject.Fqdn{}

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
	targetVirtualHostHostname tkValueObject.Fqdn,
	expectedOwnershipHash tkValueObject.Hash,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) (mappingId valueObject.MappingId, err error) {
	matchPattern, _ := valueObject.NewMappingMatchPattern("equals")
	targetType, _ := valueObject.NewMappingTargetType("inline-html")
	httpResponseCode, _ := tkValueObject.NewHttpStatusCode(200)
	targetValue, _ := valueObject.NewMappingTargetValue(
		expectedOwnershipHash.String(), targetType,
	)
	shouldUpgradeInsecureRequests := false

	return repo.mappingCmdRepo.Create(dto.NewCreateMapping(
		targetVirtualHostHostname, repo.ownershipValidationPath, matchPattern,
		targetType, &targetValue, &httpResponseCode, &shouldUpgradeInsecureRequests,
		nil, operatorAccountId, operatorIpAddress,
	))
}

func (repo *SslCmdRepo) deleteStaleOwnershipValidationMappings(
	hostname tkValueObject.Fqdn,
) error {
	existingMappingsResponse, err := repo.mappingQueryRepo.Read(dto.ReadMappingsRequest{
		Pagination:  tkDto.PaginationUnpaginated,
		Hostname:    &hostname,
		MappingPath: &repo.ownershipValidationPath,
	})
	if err != nil {
		return errors.New("ReadStaleOwnershipValidationMappingsError: " + err.Error())
	}

	for _, existingMapping := range existingMappingsResponse.Mappings {
		err := repo.mappingCmdRepo.Delete(existingMapping.Id)
		if err != nil {
			return errors.New("DeleteStaleOwnershipValidationMappingError: " + err.Error())
		}
	}
	return nil
}

func (repo *SslCmdRepo) httpFunctionalHostnamesFilter(
	vhostHostnames []tkValueObject.Fqdn,
	expectedOwnershipHash tkValueObject.Hash,
	serverPublicIpAddress tkValueObject.IpAddress,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) []tkValueObject.Fqdn {
	functionalHostnames := []tkValueObject.Fqdn{}

	serverPublicIpAddressStr := serverPublicIpAddress.String()
	expectedHashStr := expectedOwnershipHash.String()

	for _, vhostHostname := range vhostHostnames {
		vhostHostnameStr := vhostHostname.String()

		err := repo.deleteStaleOwnershipValidationMappings(vhostHostname)
		if err != nil {
			slog.Error(
				"DeleteStaleOwnershipValidationMappingsError",
				slog.String("hostname", vhostHostnameStr),
				slog.String("error", err.Error()),
			)
			continue
		}

		ownershipValidationMappingId, err := repo.createOwnershipValidationMapping(
			vhostHostname, expectedOwnershipHash, operatorAccountId, operatorIpAddress,
		)
		if err != nil {
			continue
		}

		hashUrlPath := DomainOwnershipValidationUrlPath
		hashUrlFull := "https://" + vhostHostnameStr + hashUrlPath
		curlBaseCmd := "curl -skLm 10 "
		sniFlag := "--resolve " + vhostHostnameStr + ":443:" + serverPublicIpAddressStr
		ownershipHashFound, err := tkInfra.NewShell(tkInfra.ShellSettings{
			Command:           curlBaseCmd + sniFlag + " " + hashUrlFull,
			ShouldUseSubShell: true,
		}).Run()
		if err != nil {
			hashUrlFull = "https://" + serverPublicIpAddressStr + hashUrlPath
			ownershipHashFound, err = tkInfra.NewShell(tkInfra.ShellSettings{
				Command: curlBaseCmd + "-H \"Host: " + vhostHostnameStr +
					"\" " + hashUrlFull,
				ShouldUseSubShell: true,
			}).Run()
		}

		deleteErr := repo.mappingCmdRepo.Delete(ownershipValidationMappingId)
		if deleteErr != nil {
			slog.Error(
				"DeleteOwnershipValidationMappingError",
				slog.String("error", deleteErr.Error()),
			)
		}

		if err != nil || ownershipHashFound != expectedHashStr {
			continue
		}

		functionalHostnames = append(functionalHostnames, vhostHostname)
	}

	return functionalHostnames
}

func (repo *SslCmdRepo) issueValidSsl(
	mainHostname tkValueObject.Fqdn,
	functionalHostnames []tkValueObject.Fqdn,
) error {
	mainHostnameStr := mainHostname.String()
	vhostRootDir := infraEnvs.PrimaryVirtualHostPublicDir
	if !repo.vhostHelpers.IsPrimaryVirtualHost(mainHostname) {
		vhostRootDir += "/" + mainHostnameStr
	}

	if !repo.fileClerk.FileExists(vhostRootDir) {
		return errors.New("VirtualHostRootDirNotFound")
	}

	certbotCmd := "certbot certonly --webroot --webroot-path " + vhostRootDir +
		" --agree-tos --register-unsafely-without-email --cert-name " + mainHostnameStr
	for _, functionalHostname := range functionalHostnames {
		certbotCmd += " -d " + functionalHostname.String()
	}

	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:           certbotCmd,
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return errors.New("GenerateValidSslCertError: " + err.Error())
	}

	certbotDirPath := "/etc/letsencrypt/live"
	shouldOverwrite := true

	certbotCrtFilePath := certbotDirPath + "/" + mainHostnameStr + "/fullchain.pem"
	vhostCrtFilePath := infraEnvs.PkiConfDir + "/" + mainHostnameStr + ".crt"
	err = repo.fileClerk.CreateSymlink(
		certbotCrtFilePath, vhostCrtFilePath, shouldOverwrite,
	)
	if err != nil {
		return errors.New("CreateSslCertSymlinkError: " + err.Error())
	}

	certbotKeyFilePath := certbotDirPath + "/" + mainHostnameStr + "/privkey.pem"
	vhostKeyFilePath := infraEnvs.PkiConfDir + "/" + mainHostnameStr + ".key"
	err = repo.fileClerk.CreateSymlink(
		certbotKeyFilePath, vhostKeyFilePath, shouldOverwrite,
	)
	if err != nil {
		return errors.New("CreateSslKeySymlinkError: " + err.Error())
	}

	return repo.vhostHelpers.ReloadWebServer()
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
		Pagination: tkDto.PaginationSingleItem,
		Hostname:   &createDto.VirtualHostHostname,
	})
	if err != nil {
		return sslPairId, errors.New("ReadVirtualHostEntitiesError: " + err.Error())
	}

	if len(vhostReadResponse.VirtualHosts) == 0 {
		return sslPairId, errors.New("VirtualHostNotFound")
	}

	virtualHostsHostnames := []tkValueObject.Fqdn{createDto.VirtualHostHostname}
	virtualHostsHostnames = append(
		virtualHostsHostnames, vhostReadResponse.VirtualHosts[0].AliasesHostnames...,
	)

	skipDnsOwnershipCheck := false
	envSkipDns, err := tkVoUtil.InterfaceToBool(os.Getenv("SKIP_SSL_DNS_OWNERSHIP_CHECK"))
	if err == nil && envSkipDns {
		skipDnsOwnershipCheck = true
	}

	for _, vhostHostname := range virtualHostsHostnames {
		wwwVirtualHostHostname, err := tkValueObject.NewFqdn(
			"www." + vhostHostname.String(),
		)
		if err != nil {
			slog.Debug(
				"CreatePubliclyTrustedInvalidWwwHostname",
				slog.String("fqdn", vhostHostname.String()),
			)
			continue
		}
		virtualHostsHostnames = append(
			virtualHostsHostnames, wwwVirtualHostHostname,
		)
	}

	dnsFunctionalHostnames := virtualHostsHostnames
	if !skipDnsOwnershipCheck {
		dnsFunctionalHostnames = repo.dnsFunctionalHostnamesFilter(
			virtualHostsHostnames, serverPublicIpAddress,
		)
		if len(dnsFunctionalHostnames) == 0 {
			return sslPairId, errors.New("NoSslHostnamePointingToServerIpAddress")
		}
	}

	synthesizer := tkInfra.Synthesizer{}
	dummyValue := synthesizer.PasswordFactory(32, false)
	dummyHash := infraHelper.GenStrongHash(dummyValue)

	expectedOwnershipHash, err := tkValueObject.NewHash(dummyHash)
	if err != nil {
		return sslPairId, errors.New("CreateOwnershipValidationHashError: " + err.Error())
	}

	httpFunctionalHostnames := repo.httpFunctionalHostnamesFilter(
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
		err := repo.fileClerk.UpdateFileContent(
			vhostCertFilePath, certContentStr, shouldOverwrite,
		)
		if err != nil {
			return sslPairId, errors.New("UpdateCertFileError: " + err.Error())
		}

		err = repo.fileClerk.UpdateFileContent(
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

	err = repo.vhostHelpers.ReloadWebServer()
	if err != nil {
		return sslPairId, errors.New("ReloadWebServerError: " + err.Error())
	}

	return sslPairEntity.Id, nil
}

func (repo *SslCmdRepo) ReplaceWithSelfSigned(vhostHostname tkValueObject.Fqdn) error {
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
		Pagination:     tkDto.PaginationUnpaginated,
		ParentHostname: &vhostHostname,
	})
	if err != nil {
		return errors.New("ReadAliasesError: " + err.Error())
	}

	aliasesHostnames := []tkValueObject.Fqdn{}
	for _, aliasVirtualHostEntity := range aliasesVirtualHostsReadResponse.VirtualHosts {
		aliasesHostnames = append(aliasesHostnames, aliasVirtualHostEntity.Hostname)
	}

	vhostHostnameStr := vhostHostname.String()
	vhostCertFilePath := infraEnvs.PkiConfDir + "/" + vhostHostnameStr + ".crt"
	vhostCertFileExists := repo.fileClerk.FileExists(vhostCertFilePath)
	if vhostCertFileExists {
		err := os.Remove(vhostCertFilePath)
		if err != nil {
			return errors.New("DeleteCertFileError: " + err.Error())
		}
	}

	vhostCertKeyFilePath := infraEnvs.PkiConfDir + "/" + vhostHostnameStr + ".key"
	vhostCertKeyFileExists := repo.fileClerk.FileExists(vhostCertKeyFilePath)
	if vhostCertKeyFileExists {
		err := os.Remove(vhostCertKeyFilePath)
		if err != nil {
			return errors.New("DeleteCertKeyFileError: " + err.Error())
		}
	}

	pkiConfDir, err := tkValueObject.NewUnixAbsoluteFilePath(infraEnvs.PkiConfDir, false)
	if err != nil {
		return errors.New("PkiConfDirError: " + err.Error())
	}

	err = infraHelper.CreateSelfSignedSsl(pkiConfDir, vhostHostname, aliasesHostnames)
	if err != nil {
		return errors.New("CreateSelfSignedSslError: " + err.Error())
	}

	return repo.vhostHelpers.ReloadWebServer()
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
