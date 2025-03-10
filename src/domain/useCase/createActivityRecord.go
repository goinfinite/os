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

func (uc *CreateSecurityActivityRecord) UpdateAccount(updateDto dto.UpdateAccount) {
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewAccountSri(updateDto.AccountId),
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
	operatorAccountId := createDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("CronCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewCronSri(operatorAccountId, cronId),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateCron(updateDto dto.UpdateCron) {
	operatorAccountId := updateDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("CronUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewCronSri(operatorAccountId, updateDto.Id),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteCron(deleteDto dto.DeleteCron) {
	operatorAccountId := deleteDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("CronDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewCronSri(operatorAccountId, *deleteDto.Id),
		},
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateDatabase(createDto dto.CreateDatabase) {
	operatorAccountId := createDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("DatabaseCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(operatorAccountId, createDto.DatabaseName),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteDatabase(deleteDto dto.DeleteDatabase) {
	operatorAccountId := deleteDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("DatabaseDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(operatorAccountId, deleteDto.DatabaseName),
		},
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateDatabaseUser(
	createDto dto.CreateDatabaseUser,
) {
	operatorAccountId := createDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("DatabaseUserCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(operatorAccountId, createDto.DatabaseName),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteDatabaseUser(
	deleteDto dto.DeleteDatabaseUser,
) {
	operatorAccountId := deleteDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("DatabaseUserDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(operatorAccountId, deleteDto.DatabaseName),
			valueObject.NewDatabaseUserSri(operatorAccountId, deleteDto.Username),
		},
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) InstallMarketplaceCatalogItem(
	installDto dto.InstallMarketplaceCatalogItem,
) {
	operatorAccountId := installDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("MarketplaceCatalogItemInstalled")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMarketplaceCatalogItemSri(
				operatorAccountId, installDto.Id, installDto.Slug,
			),
		},
		RecordDetails:     installDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &installDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMarketplaceInstalledItem(
	installDto dto.DeleteMarketplaceInstalledItem,
) {
	operatorAccountId := installDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("MarketplaceInstalledItemDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMarketplaceInstalledItemSri(
				operatorAccountId, installDto.InstalledId,
			),
		},
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &installDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdatePhpConfigs(
	updateDto dto.UpdatePhpConfigs,
	configType string,
) {
	operatorAccountId := updateDto.OperatorAccountId

	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewPhpRuntimeSri(operatorAccountId, updateDto.Hostname),
		},
		OperatorAccountId: &operatorAccountId,
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
	operatorAccountId := createDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("InstallableServiceCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewInstallableServiceSri(operatorAccountId, createDto.Name),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateCustomService(
	createDto dto.CreateCustomService,
) {
	operatorAccountId := createDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("CustomServiceCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewCustomServiceSri(operatorAccountId, createDto.Name),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateService(updateDto dto.UpdateService) {
	operatorAccountId := updateDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("ServiceUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewInstalledServiceSri(operatorAccountId, updateDto.Name),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteService(deleteDto dto.DeleteService) {
	operatorAccountId := deleteDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("ServiceDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewInstalledServiceSri(operatorAccountId, deleteDto.Name),
		},
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateSslPair(
	createDto dto.CreateSslPair,
	sslPairId valueObject.SslPairId,
) {
	operatorAccountId := createDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("SslPairCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewSslSri(operatorAccountId, sslPairId),
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
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteSslPair(deleteDto dto.DeleteSslPair) {
	operatorAccountId := deleteDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("SslPairDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewSslSri(operatorAccountId, deleteDto.SslPairId),
		},
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteSslPairVhosts(
	deleteDto dto.DeleteSslPairVhosts,
) {
	operatorAccountId := deleteDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("SslPairVhostsDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewSslSri(operatorAccountId, deleteDto.SslPairId),
		},
		RecordDetails: map[string]interface{}{
			"sslPairVhosts": deleteDto.VirtualHostsHostnames,
		},
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateVirtualHost(
	createDto dto.CreateVirtualHost,
) {
	operatorAccountId := createDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("VirtualHostCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewVirtualHostSri(operatorAccountId, createDto.Hostname),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteVirtualHost(
	deleteDto dto.DeleteVirtualHost,
) {
	operatorAccountId := deleteDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("VirtualHostDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewVirtualHostSri(operatorAccountId, deleteDto.Hostname),
		},
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateMapping(
	createDto dto.CreateMapping,
	mappingId valueObject.MappingId,
) {
	operatorAccountId := createDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("MappingCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(operatorAccountId, mappingId),
		},
		RecordDetails:     createDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMapping(deleteDto dto.DeleteMapping) {
	operatorAccountId := deleteDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("MappingDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(operatorAccountId, deleteDto.MappingId),
		},
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateUnixFile(createDto dto.CreateUnixFile) {
	operatorAccountId := createDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("UnixFileCreated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{},
		RecordDetails:     createDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteUnixFiles(deleteDto dto.DeleteUnixFiles) {
	operatorAccountId := deleteDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("UnixFileDeleted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{},
		RecordDetails:     deleteDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateUnixFiles(updateDto dto.UpdateUnixFiles) {
	operatorAccountId := updateDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("UnixFileUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{},
		RecordDetails:     updateDto,
		OperatorAccountId: &operatorAccountId,
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
	operatorAccountId := copyDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("UnixFileCopied")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{},
		RecordDetails:     copyDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &copyDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CompressUnixFile(
	compressDto dto.CompressUnixFiles,
) {
	operatorAccountId := compressDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("UnixFilesCompressed")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{},
		RecordDetails:     compressDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &compressDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) ExtractUnixFile(
	extractDto dto.ExtractUnixFiles,
) {
	operatorAccountId := extractDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("UnixFileExtracted")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		RecordDetails:     extractDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &extractDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UploadUnixFiles(
	uploadDto dto.UploadUnixFiles,
) {
	operatorAccountId := uploadDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("UnixFilesUploaded")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: &operatorAccountId,
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
