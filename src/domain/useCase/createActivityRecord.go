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
			slog.Any("error", err),
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
	updateDto dto.UpdateAccount,
) {
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

func (uc *CreateSecurityActivityRecord) DeleteAccount(
	deleteDto dto.DeleteAccount,
) {
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

func (uc *CreateSecurityActivityRecord) UpdateCron(
	updateDto dto.UpdateCron,
) {
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

func (uc *CreateSecurityActivityRecord) DeleteCron(
	deleteDto dto.DeleteCron,
) {
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

func (uc *CreateSecurityActivityRecord) CreateDatabase(
	createDto dto.CreateDatabase,
) {
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

func (uc *CreateSecurityActivityRecord) DeleteDatabase(
	deleteDto dto.DeleteDatabase,
) {
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

type phpConfigNameEnum string

const (
	Version  phpConfigNameEnum = "version"
	Modules  phpConfigNameEnum = "modules"
	Settings phpConfigNameEnum = "settings"
)

func (uc *CreateSecurityActivityRecord) UpdatePhpConfigs(
	updateDto dto.UpdatePhpConfigs,
	configName phpConfigNameEnum,
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
	var details interface{} = map[string]interface{}{
		"version": updateDto.PhpVersion.String(),
	}

	switch configName {
	case Modules:
		codeStr = "PhpModulesUpdated"
		details = updateDto.PhpModules
	case Settings:
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

func (uc *CreateSecurityActivityRecord) UpdateService(
	updateDto dto.UpdateService,
) {
	operatorAccountId := updateDto.OperatorAccountId

	recordCode, _ := valueObject.NewActivityRecordCode("ServiceUpdated")
	createRecordDto := dto.CreateActivityRecord{
		RecordLevel: uc.recordLevel,
		RecordCode:  recordCode,
		AffectedResources: []valueObject.SystemResourceIdentifier{
			valueObject.NewInstallableServiceSri(operatorAccountId, updateDto.Name),
		},
		RecordDetails:     updateDto,
		OperatorAccountId: &operatorAccountId,
		OperatorIpAddress: &updateDto.OperatorIpAddress,
	}

	uc.createActivityRecord(createRecordDto)
}
