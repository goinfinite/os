package useCase

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type CreateSecurityActivityRecord struct {
	activityRecordCmdRepo repository.ActivityRecordCmdRepo
	recordLevel           valueObject.ActivityRecordLevel
}

func NewCreateSecurityActivityRecord(
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
) *CreateSecurityActivityRecord {
	recordLevel, _ := valueObject.NewActivityRecordLevel("SEC")
	return &CreateSecurityActivityRecord{
		activityRecordCmdRepo: activityRecordCmdRepo,
		recordLevel:           recordLevel,
	}
}

func (uc *CreateSecurityActivityRecord) createActivityRecord(
	createDto dto.CreateActivityRecord,
) {
	err := uc.activityRecordCmdRepo.Create(createDto)
	if err != nil {
		slog.Debug(
			"CreateSecurityActivityRecordError",
			slog.Any("createDto", createDto),
			slog.String("err", err.Error()),
		)
	}
}

func (uc *CreateSecurityActivityRecord) CreateSessionToken(
	recordCode valueObject.ActivityRecordCode,
	createDto dto.CreateSessionToken,
) {
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		RecordDetails:     createDto.Username,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateAccount(
	createDto dto.CreateAccount,
	accountId valueObject.AccountId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("AccountCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewAccountSri(accountId),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateAccount(
	accountId valueObject.AccountId,
	updateDto dto.UpdateAccount,
) {
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewAccountSri(accountId),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	codeStr := "AccountUpdated"
	if updateDto.Password != nil {
		codeStr = "AccountPasswordUpdated"
		createRecordDto.RecordDetails = nil
	}

	if updateDto.ShouldUpdateApiKey != nil && *updateDto.ShouldUpdateApiKey {
		codeStr = "AccountApiKeyUpdated"
		createRecordDto.RecordDetails = nil
	}

	recordCode, _ := valueObject.NewActivityRecordCode(codeStr)
	createRecordDto.RecordCode = recordCode

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteAccount(deleteDto dto.DeleteAccount) {
	recordCode, _ := valueObject.NewActivityRecordCode("AccountDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewAccountSri(deleteDto.AccountId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateSecureAccessPublicKey(
	createDto dto.CreateSecureAccessPublicKey,
	keyId valueObject.SecureAccessPublicKeyId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("SecureAccessPublicKeyCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewSecureAccessPublicKeySri(createDto.AccountId, keyId),
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteSecureAccessPublicKey(
	deleteDto dto.DeleteSecureAccessPublicKey,
	accountId valueObject.AccountId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("SecureAccessPublicKeyDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewSecureAccessPublicKeySri(accountId, deleteDto.Id),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateCron(
	createDto dto.CreateCron,
	cronId valueObject.CronId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("CronCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewCronSri(createDto.OperatorAccountId, cronId),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateCron(updateDto dto.UpdateCron) {
	recordCode, _ := valueObject.NewActivityRecordCode("CronUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewCronSri(updateDto.OperatorAccountId, updateDto.Id),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteCron(deleteDto dto.DeleteCron) {
	recordCode, _ := valueObject.NewActivityRecordCode("CronDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewCronSri(deleteDto.OperatorAccountId, *deleteDto.Id),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateDatabase(createDto dto.CreateDatabase) {
	recordCode, _ := valueObject.NewActivityRecordCode("DatabaseCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(createDto.OperatorAccountId, createDto.DatabaseName),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteDatabase(deleteDto dto.DeleteDatabase) {
	recordCode, _ := valueObject.NewActivityRecordCode("DatabaseDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(deleteDto.OperatorAccountId, deleteDto.DatabaseName),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateDatabaseUser(
	createDto dto.CreateDatabaseUser,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("DatabaseUserCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(createDto.OperatorAccountId, createDto.DatabaseName),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteDatabaseUser(
	deleteDto dto.DeleteDatabaseUser,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("DatabaseUserDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(deleteDto.OperatorAccountId, deleteDto.DatabaseName),
			valueObject.NewDatabaseUserSri(deleteDto.OperatorAccountId, deleteDto.Username),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) InstallMarketplaceCatalogItem(
	installDto dto.InstallMarketplaceCatalogItem,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MarketplaceCatalogItemInstalled")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMarketplaceCatalogItemSri(
				installDto.OperatorAccountId, installDto.Id, installDto.Slug,
			),
		},
		RecordDetails:     installDto,
		OperatorAccountId: &installDto.OperatorAccountId,
		OperatorIpAddress: &installDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMarketplaceInstalledItem(
	installDto dto.DeleteMarketplaceInstalledItem,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MarketplaceInstalledItemDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMarketplaceInstalledItemSri(
				installDto.OperatorAccountId, installDto.InstalledId,
			),
		},
		OperatorAccountId: &installDto.OperatorAccountId,
		OperatorIpAddress: &installDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdatePhpConfigs(
	updateDto dto.UpdatePhpConfigs,
	configType string,
) {
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewPhpRuntimeSri(updateDto.OperatorAccountId, updateDto.Hostname),
		},
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	codeStr := "PhpVersionUpdated"
	var details interface{} = map[string]string{
		"version": updateDto.PhpVersion.String(),
	}

	switch configType {
	case "modules":
		codeStr = "PhpModulesUpdated"
		details = updateDto.PhpModules
	case "settings":
		codeStr = "PhpSettingsUpdated"
		details = updateDto.PhpSettings
	}

	recordCode, _ := valueObject.NewActivityRecordCode(codeStr)
	createRecordDto.RecordCode = recordCode
	createRecordDto.RecordDetails = details

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateInstallableService(
	createDto dto.CreateInstallableService,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("InstallableServiceCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewInstallableServiceSri(createDto.OperatorAccountId, createDto.Name),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateCustomService(
	createDto dto.CreateCustomService,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("CustomServiceCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewCustomServiceSri(createDto.OperatorAccountId, createDto.Name),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateService(updateDto dto.UpdateService) {
	recordCode, _ := valueObject.NewActivityRecordCode("ServiceUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewInstalledServiceSri(updateDto.OperatorAccountId, updateDto.Name),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteService(deleteDto dto.DeleteService) {
	recordCode, _ := valueObject.NewActivityRecordCode("ServiceDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewInstalledServiceSri(deleteDto.OperatorAccountId, deleteDto.Name),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateSslPair(
	createDto dto.CreateSslPair,
	sslPairId valueObject.SslPairId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("SslPairCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewSslSri(createDto.OperatorAccountId, sslPairId),
		},
		RecordDetails: map[string]interface{}{
			"virtualHostHostnames": createDto.VirtualHostsHostnames,
			"certificate": map[string]interface{}{
				"commonName": createDto.Certificate.CommonName,
				"altNames":   createDto.Certificate.AltNames,
				"authority":  createDto.Certificate.CertificateAuthority,
				"issuedAt":   createDto.Certificate.IssuedAt,
				"expiresAt":  createDto.Certificate.ExpiresAt,
			},
		},
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreatePubliclyTrustedSslPair(
	createDto dto.CreatePubliclyTrustedSslPair,
	sslPairId valueObject.SslPairId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("PubliclyTrustedSslPairCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewSslSri(createDto.OperatorAccountId, sslPairId),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteSslPair(deleteDto dto.DeleteSslPair) {
	recordCode, _ := valueObject.NewActivityRecordCode("SslPairDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewSslSri(deleteDto.OperatorAccountId, deleteDto.SslPairId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteSslPairVhosts(
	deleteDto dto.DeleteSslPairVhosts,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("SslPairVhostsDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewSslSri(deleteDto.OperatorAccountId, deleteDto.SslPairId),
		},
		RecordDetails: map[string]interface{}{
			"sslPairVhosts": deleteDto.VirtualHostsHostnames,
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateVirtualHost(
	createDto dto.CreateVirtualHost,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("VirtualHostCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewVirtualHostSri(createDto.OperatorAccountId, createDto.Hostname),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateVirtualHost(
	updateDto dto.UpdateVirtualHost,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("VirtualHostUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewVirtualHostSri(updateDto.OperatorAccountId, updateDto.Hostname),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteVirtualHost(
	deleteDto dto.DeleteVirtualHost,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("VirtualHostDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewVirtualHostSri(deleteDto.OperatorAccountId, deleteDto.Hostname),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateMapping(
	createDto dto.CreateMapping,
	mappingId valueObject.MappingId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(createDto.OperatorAccountId, mappingId),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateMapping(updateDto dto.UpdateMapping) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(updateDto.OperatorAccountId, updateDto.Id),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMapping(deleteDto dto.DeleteMapping) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(deleteDto.OperatorAccountId, deleteDto.MappingId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateUnixFile(createDto dto.CreateUnixFile) {
	recordCode, _ := valueObject.NewActivityRecordCode("UnixFileCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteUnixFiles(deleteDto dto.DeleteUnixFiles) {
	recordCode, _ := valueObject.NewActivityRecordCode("UnixFileDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{},
		RecordDetails:     deleteDto,
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateUnixFiles(updateDto dto.UpdateUnixFiles) {
	recordCode, _ := valueObject.NewActivityRecordCode("UnixFileUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	details := map[string]interface{}{
		"sourcePaths": updateDto.SourcePaths,
	}
	if updateDto.DestinationPath != nil {
		details["destinationPath"] = updateDto.DestinationPath
	}
	if updateDto.Permissions != nil {
		details["permissions"] = updateDto.Permissions
	}
	if updateDto.EncodedContent != nil {
		details["contentWasUpdated"] = true
	}
	createRecordDto.RecordDetails = details

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CopyUnixFile(copyDto dto.CopyUnixFile) {
	recordCode, _ := valueObject.NewActivityRecordCode("UnixFileCopied")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{},
		RecordDetails:     copyDto,
		OperatorAccountId: &copyDto.OperatorAccountId,
		OperatorIpAddress: &copyDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CompressUnixFile(
	compressDto dto.CompressUnixFiles,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("UnixFilesCompressed")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{},
		RecordDetails:     compressDto,
		OperatorAccountId: &compressDto.OperatorAccountId,
		OperatorIpAddress: &compressDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) ExtractUnixFile(
	extractDto dto.ExtractUnixFiles,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("UnixFileExtracted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		RecordDetails:     extractDto,
		OperatorAccountId: &extractDto.OperatorAccountId,
		OperatorIpAddress: &extractDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UploadUnixFiles(
	uploadDto dto.UploadUnixFiles,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("UnixFilesUploaded")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &uploadDto.OperatorAccountId,
		OperatorIpAddress: &uploadDto.OperatorIpAddress,
	}

	details := map[string]interface{}{"destinationPath": uploadDto.DestinationPath}

	fileNames := []valueObject.UnixFileName{}
	for _, fileStreamHandler := range uploadDto.FileStreamHandlers {
		fileNames = append(fileNames, fileStreamHandler.Name)
	}
	details["fileNames"] = fileNames

	createRecordDto.RecordDetails = details

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateMappingSecurityRule(
	createDto dto.CreateMappingSecurityRule,
	mappingSecurityRuleId valueObject.MappingSecurityRuleId,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingSecurityRuleCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSecurityRuleSri(createDto.OperatorAccountId, mappingSecurityRuleId),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &createDto.OperatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateMappingSecurityRule(
	updateDto dto.UpdateMappingSecurityRule,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingSecurityRuleUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSecurityRuleSri(updateDto.OperatorAccountId, updateDto.Id),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &updateDto.OperatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMappingSecurityRule(
	deleteDto dto.DeleteMappingSecurityRule,
) {
	recordCode, _ := valueObject.NewActivityRecordCode("MappingSecurityRuleDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSecurityRuleSri(deleteDto.OperatorAccountId, deleteDto.SecurityRuleId),
		},
		OperatorAccountId: &deleteDto.OperatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}
