package useCase

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateSecurityActivityRecord struct {
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo
	recordLevel           tkValueObject.ActivityRecordLevel
}

func NewCreateSecurityActivityRecord(
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
) *CreateSecurityActivityRecord {
	recordLevel := tkValueObject.ActivityRecordLevelSecurity
	return &CreateSecurityActivityRecord{
		activityRecordCmdRepo: activityRecordCmdRepo,
		recordLevel:           recordLevel,
	}
}

func (uc *CreateSecurityActivityRecord) createActivityRecord(
	createDto tkDto.CreateActivityRecord,
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
	recordCode tkValueObject.ActivityRecordCode,
	createDto dto.CreateSessionToken,
) {
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		RecordDetails:     createDto.Username,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateAccount(
	createDto dto.CreateAccount,
	accountId tkValueObject.AccountId,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("AccountCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			tkValueObject.NewSriAccount(accountId),
		},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateAccount(
	accountId tkValueObject.AccountId,
	updateDto dto.UpdateAccount,
) {
	operatorSri := tkValueObject.NewSriAccount(updateDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			tkValueObject.NewSriAccount(accountId),
		},
		RecordDetails:     updateDto,
		OperatorSri:       &operatorSri,
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

	recordCode, _ := tkValueObject.NewActivityRecordCode(codeStr)
	createRecordDto.RecordCode = recordCode

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteAccount(deleteDto dto.DeleteAccount) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("AccountDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			tkValueObject.NewSriAccount(deleteDto.AccountId),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateSecureAccessPublicKey(
	createDto dto.CreateSecureAccessPublicKey,
	keyId valueObject.SecureAccessPublicKeyId,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("SecureAccessPublicKeyCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewSecureAccessPublicKeySri(createDto.AccountId, keyId),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteSecureAccessPublicKey(
	deleteDto dto.DeleteSecureAccessPublicKey,
	accountId tkValueObject.AccountId,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("SecureAccessPublicKeyDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewSecureAccessPublicKeySri(accountId, deleteDto.Id),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateCron(
	createDto dto.CreateCron,
	cronId valueObject.CronId,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("CronCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewCronSri(createDto.OperatorAccountId, cronId),
		},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateCron(updateDto dto.UpdateCron) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("CronUpdated")
	operatorSri := tkValueObject.NewSriAccount(updateDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewCronSri(updateDto.OperatorAccountId, updateDto.Id),
		},
		RecordDetails:     updateDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteCron(deleteDto dto.DeleteCron) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("CronDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewCronSri(deleteDto.OperatorAccountId, *deleteDto.Id),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateDatabase(createDto dto.CreateDatabase) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("DatabaseCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(createDto.OperatorAccountId, createDto.DatabaseName),
		},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteDatabase(deleteDto dto.DeleteDatabase) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("DatabaseDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(deleteDto.OperatorAccountId, deleteDto.DatabaseName),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateDatabaseUser(
	createDto dto.CreateDatabaseUser,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("DatabaseUserCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(createDto.OperatorAccountId, createDto.DatabaseName),
		},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteDatabaseUser(
	deleteDto dto.DeleteDatabaseUser,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("DatabaseUserDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewDatabaseSri(deleteDto.OperatorAccountId, deleteDto.DatabaseName),
			valueObject.NewDatabaseUserSri(deleteDto.OperatorAccountId, deleteDto.Username),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) InstallMarketplaceCatalogItem(
	installDto dto.InstallMarketplaceCatalogItem,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("MarketplaceCatalogItemInstalled")
	operatorSri := tkValueObject.NewSriAccount(installDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewMarketplaceCatalogItemSri(
				installDto.OperatorAccountId, installDto.Id, installDto.Slug,
			),
		},
		RecordDetails:     installDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &installDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMarketplaceInstalledItem(
	installDto dto.DeleteMarketplaceInstalledItem,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("MarketplaceInstalledItemDeleted")
	operatorSri := tkValueObject.NewSriAccount(installDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewMarketplaceInstalledItemSri(
				installDto.OperatorAccountId, installDto.InstalledId,
			),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &installDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdatePhpConfigs(
	updateDto dto.UpdatePhpConfigs,
	configType string,
) {
	operatorSri := tkValueObject.NewSriAccount(updateDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewPhpRuntimeSri(updateDto.OperatorAccountId, updateDto.Hostname),
		},
		OperatorSri:       &operatorSri,
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

	recordCode, _ := tkValueObject.NewActivityRecordCode(codeStr)
	createRecordDto.RecordCode = recordCode
	createRecordDto.RecordDetails = details

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateInstallableService(
	createDto dto.CreateInstallableService,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("InstallableServiceCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewInstallableServiceSri(createDto.OperatorAccountId, createDto.Name),
		},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateCustomService(
	createDto dto.CreateCustomService,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("CustomServiceCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewCustomServiceSri(createDto.OperatorAccountId, createDto.Name),
		},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateService(updateDto dto.UpdateService) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("ServiceUpdated")
	operatorSri := tkValueObject.NewSriAccount(updateDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewInstalledServiceSri(updateDto.OperatorAccountId, updateDto.Name),
		},
		RecordDetails:     updateDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteService(deleteDto dto.DeleteService) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("ServiceDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewInstalledServiceSri(deleteDto.OperatorAccountId, deleteDto.Name),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateSslPair(
	createDto dto.CreateSslPair,
	sslPairId valueObject.SslPairId,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("SslPairCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
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
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreatePubliclyTrustedSslPair(
	createDto dto.CreatePubliclyTrustedSslPair,
	sslPairId valueObject.SslPairId,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("PubliclyTrustedSslPairCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewSslSri(createDto.OperatorAccountId, sslPairId),
		},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteSslPair(deleteDto dto.DeleteSslPair) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("SslPairDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewSslSri(deleteDto.OperatorAccountId, deleteDto.SslPairId),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteSslPairVhosts(
	deleteDto dto.DeleteSslPairVhosts,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("SslPairVhostsDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewSslSri(deleteDto.OperatorAccountId, deleteDto.SslPairId),
		},
		RecordDetails: map[string]interface{}{
			"sslPairVhosts": deleteDto.VirtualHostsHostnames,
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateVirtualHost(
	createDto dto.CreateVirtualHost,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("VirtualHostCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewVirtualHostSri(createDto.OperatorAccountId, createDto.Hostname),
		},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateVirtualHost(
	updateDto dto.UpdateVirtualHost,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("VirtualHostUpdated")
	operatorSri := tkValueObject.NewSriAccount(updateDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewVirtualHostSri(updateDto.OperatorAccountId, updateDto.Hostname),
		},
		RecordDetails:     updateDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteVirtualHost(
	deleteDto dto.DeleteVirtualHost,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("VirtualHostDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewVirtualHostSri(deleteDto.OperatorAccountId, deleteDto.Hostname),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateMapping(
	createDto dto.CreateMapping,
	mappingId valueObject.MappingId,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("MappingCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(createDto.OperatorAccountId, mappingId),
		},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateMapping(updateDto dto.UpdateMapping) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("MappingUpdated")
	operatorSri := tkValueObject.NewSriAccount(updateDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(updateDto.OperatorAccountId, updateDto.Id),
		},
		RecordDetails:     updateDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMapping(deleteDto dto.DeleteMapping) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("MappingDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewMappingSri(deleteDto.OperatorAccountId, deleteDto.MappingId),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CreateUnixFile(createDto dto.CreateUnixFile) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("UnixFileCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteUnixFiles(deleteDto dto.DeleteUnixFiles) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("UnixFileDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{},
		RecordDetails:     deleteDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateUnixFiles(updateDto dto.UpdateUnixFiles) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("UnixFileUpdated")
	operatorSri := tkValueObject.NewSriAccount(updateDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{},
		RecordDetails:     updateDto,
		OperatorSri:       &operatorSri,
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
	recordCode, _ := tkValueObject.NewActivityRecordCode("UnixFileCopied")
	operatorSri := tkValueObject.NewSriAccount(copyDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{},
		RecordDetails:     copyDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &copyDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) CompressUnixFile(
	compressDto dto.CompressUnixFiles,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("UnixFilesCompressed")
	operatorSri := tkValueObject.NewSriAccount(compressDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{},
		RecordDetails:     compressDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &compressDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) ExtractUnixFile(
	extractDto dto.ExtractUnixFiles,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("UnixFileExtracted")
	operatorSri := tkValueObject.NewSriAccount(extractDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		RecordDetails:     extractDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &extractDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UploadUnixFiles(
	uploadDto dto.UploadUnixFiles,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("UnixFilesUploaded")
	operatorSri := tkValueObject.NewSriAccount(uploadDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel:       uc.recordLevel,
		RecordCode:        recordCode,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &uploadDto.OperatorIpAddress,
	}

	details := map[string]interface{}{"destinationPath": uploadDto.DestinationPath}

	fileNames := []tkValueObject.UnixFileName{}
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
	recordCode, _ := tkValueObject.NewActivityRecordCode("MappingSecurityRuleCreated")
	operatorSri := tkValueObject.NewSriAccount(createDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewMappingSecurityRuleSri(createDto.OperatorAccountId, mappingSecurityRuleId),
		},
		RecordDetails:     createDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &createDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) UpdateMappingSecurityRule(
	updateDto dto.UpdateMappingSecurityRule,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("MappingSecurityRuleUpdated")
	operatorSri := tkValueObject.NewSriAccount(updateDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewMappingSecurityRuleSri(updateDto.OperatorAccountId, updateDto.Id),
		},
		RecordDetails:     updateDto,
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}

func (uc *CreateSecurityActivityRecord) DeleteMappingSecurityRule(
	deleteDto dto.DeleteMappingSecurityRule,
) {
	recordCode, _ := tkValueObject.NewActivityRecordCode("MappingSecurityRuleDeleted")
	operatorSri := tkValueObject.NewSriAccount(deleteDto.OperatorAccountId)
	createRecordDto := tkDto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []tkValueObject.SystemResourceIdentifier{
			valueObject.NewMappingSecurityRuleSri(deleteDto.OperatorAccountId, deleteDto.SecurityRuleId),
		},
		OperatorSri:       &operatorSri,
		OperatorIpAddress: &deleteDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}
