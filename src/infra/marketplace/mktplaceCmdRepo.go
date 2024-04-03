package mktplaceInfra

import (
	"errors"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

type MktplaceCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	queryRepo       *MktplaceQueryRepo
}

func NewMktplaceCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MktplaceCmdRepo {
	mktplaceQueryRepo := NewMktplaceQueryRepo(persistentDbSvc)

	return &MktplaceCmdRepo{
		persistentDbSvc: persistentDbSvc,
		queryRepo:       mktplaceQueryRepo,
	}
}

func (repo *MktplaceCmdRepo) getDataFieldsAsMap(
	dataFields []valueObject.DataField,
) map[string]string {
	dataFieldMap := map[string]string{}

	for _, dataField := range dataFields {
		dataFieldMap[dataField.Key.String()] = dataField.Value.String()
	}

	return dataFieldMap
}

func (repo *MktplaceCmdRepo) InstallItem(
	installMktplaceCatalogItem dto.InstallMarketplaceCatalogItem,
) error {
	mktplaceCatalogItem, err := repo.queryRepo.GetItemById(
		installMktplaceCatalogItem.Id,
	)
	if err != nil {
		return errors.New("MktplaceCatalogItemNotFound")
	}

	for _, requiredSvcName := range mktplaceCatalogItem.Services {
		requiredSvcAutoCreateMapping := false
		requiredService := dto.NewCreateInstallableService(
			requiredSvcName,
			nil,
			nil,
			nil,
			requiredSvcAutoCreateMapping,
		)

		err := servicesInfra.CreateInstallable(requiredService)
		if err != nil {
			return errors.New("InstallRequiredService: " + err.Error())
		}
	}

	// Criar os mappings que o mktplaceCatalogItem exige.

	dataFieldsMap := repo.getDataFieldsAsMap(installMktplaceCatalogItem.DataFields)
	for _, cmdStep := range mktplaceCatalogItem.CmdSteps {
		cmdStepStr := cmdStep.String()
		cmdStepRequiredDataField, err := infraHelper.GetRegexFirstGroup(
			cmdStepStr,
			`%(.+)%`,
		)
		if err == nil {
			requiredDataFieldValue := dataFieldsMap[cmdStepRequiredDataField]
			cmdStepWithDataField := strings.ReplaceAll(
				cmdStepStr,
				"%"+cmdStepRequiredDataField+"%",
				requiredDataFieldValue,
			)
			cmdStepStr = cmdStepWithDataField
		}

		_, err = infraHelper.RunCmdWithSubShell(cmdStepStr)
		if err != nil {
			return errors.New("RunCmdStepError: " + err.Error())
		}
	}

	nowUnixTime := time.Now().Unix()
	createdAt := valueObject.UnixTime(nowUnixTime)
	updatedAt := valueObject.UnixTime(nowUnixTime)

	mktplaceInstalledItem := entity.NewMarketplaceInstalledItem(
		mktplaceCatalogItem.Id,
		mktplaceCatalogItem.Name,
		mktplaceCatalogItem.Type,
		installMktplaceCatalogItem.RootDirectory,
		mktplaceCatalogItem.Services,
		[]entity.Mapping{},
		mktplaceCatalogItem.AvatarUrl,
		createdAt,
		updatedAt,
	)

	mktplaceInstalledItemModel, err := dbModel.MarketplaceInstalledItem{}.ToModel(
		mktplaceInstalledItem,
	)
	if err != nil {
		return err
	}

	err = repo.persistentDbSvc.Handler.Create(&mktplaceInstalledItemModel).Error
	if err != nil {
		return err
	}
	return nil
}
