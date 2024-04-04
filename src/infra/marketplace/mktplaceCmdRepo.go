package marketplaceInfra

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	filesInfra "github.com/speedianet/os/src/infra/files"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

type MarketplaceCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	queryRepo       *MarketplaceQueryRepo
}

func NewMarketplaceCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceCmdRepo {
	marketplaceQueryRepo := NewMarketplaceQueryRepo(persistentDbSvc)

	return &MarketplaceCmdRepo{
		persistentDbSvc: persistentDbSvc,
		queryRepo:       marketplaceQueryRepo,
	}
}

func (repo *MarketplaceCmdRepo) getDataFieldsAsMap(
	dataFields []valueObject.MarketplaceItemDataField,
) map[string]string {
	dataFieldMap := map[string]string{}

	for _, dataField := range dataFields {
		dataFieldMap[dataField.Key.String()] = dataField.Value.String()
	}

	return dataFieldMap
}

func (repo *MarketplaceCmdRepo) getCmdStepWithDataFields(
	cmdStep valueObject.MarketplaceItemInstallStep,
	dataFieldsMap map[string]string,
) (string, error) {
	cmdStepStr := cmdStep.String()
	cmdStepRequiredDataFields, _ := infraHelper.GetRegexFirstGroup(
		cmdStepStr,
		`%(.*?)%`,
	)

	cmdStepWithDataField := cmdStepStr
	for _, cmdStepRequiredDataField := range cmdStepRequiredDataFields {
		requiredDataFieldValue := dataFieldsMap[cmdStepRequiredDataField]
		cmdStepWithDataField = strings.ReplaceAll(
			cmdStepWithDataField,
			"%"+cmdStepRequiredDataField+"%",
			requiredDataFieldValue,
		)
	}

	return cmdStepWithDataField, nil
}

func (repo *MarketplaceCmdRepo) moveMarketplaceItemDir(
	rootDirectory valueObject.UnixFilePath,
	marketplaceItemName valueObject.MarketplaceItemName,
) error {
	marketplaceItemSrcPath, _ := valueObject.NewUnixFilePath("/speedia/" + marketplaceItemName.String())
	marketplaceItemDestinationPath, _ := valueObject.NewUnixFilePath(rootDirectory.String())

	filesCmdRepo := filesInfra.FilesCmdRepo{}
	shouldOverwrite := true
	return filesCmdRepo.Move(marketplaceItemSrcPath, marketplaceItemDestinationPath, shouldOverwrite)
}

func (repo *MarketplaceCmdRepo) InstallItem(
	installMarketplaceCatalogItem dto.InstallMarketplaceCatalogItem,
) error {
	marketplaceCatalogItem, err := repo.queryRepo.GetItemById(
		installMarketplaceCatalogItem.Id,
	)
	if err != nil {
		return errors.New("MarketplaceCatalogItemNotFound")
	}

	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
	for _, requiredSvcName := range marketplaceCatalogItem.ServiceNames {
		_, err := servicesQueryRepo.GetByName(requiredSvcName)
		if err == nil {
			continue
		}

		requiredSvcAutoCreateMapping := false
		requiredService := dto.NewCreateInstallableService(
			requiredSvcName,
			nil,
			nil,
			nil,
			requiredSvcAutoCreateMapping,
		)

		err = servicesCmdRepo.CreateInstallable(requiredService)
		if err != nil {
			return errors.New("InstallRequiredServiceError: " + err.Error())
		}
	}

	dataFieldsMap := repo.getDataFieldsAsMap(installMarketplaceCatalogItem.DataFields)
	for _, cmdStep := range marketplaceCatalogItem.CmdSteps {
		cmdStepRequiredDataFields, err := repo.getCmdStepWithDataFields(
			cmdStep,
			dataFieldsMap,
		)
		if err != nil {
			return errors.New("GetCmdStepWithDataFieldsError: " + err.Error())
		}

		_, err = infraHelper.RunCmdWithSubShell(cmdStepRequiredDataFields)
		if err != nil {
			return errors.New(
				"RunCmdStepError (" + cmdStepRequiredDataFields + "): " + err.Error(),
			)
		}
	}

	err = repo.moveMarketplaceItemDir(
		installMarketplaceCatalogItem.RootDirectory,
		marketplaceCatalogItem.Name,
	)
	if err != nil {
		return err
	}

	for _, marketplaceItemMapping := range marketplaceCatalogItem.Mappings {
		createMarketplaceItemMapping := dto.NewCreateMapping(
			installMarketplaceCatalogItem.Hostname,
			marketplaceItemMapping.Path,
			marketplaceItemMapping.MatchPattern,
			marketplaceItemMapping.TargetType,
			marketplaceItemMapping.TargetServiceName,
			marketplaceItemMapping.TargetUrl,
			marketplaceItemMapping.TargetHttpResponseCode,
			marketplaceItemMapping.TargetInlineHtmlContent,
		)

		vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}
		err = vhostCmdRepo.CreateMapping(createMarketplaceItemMapping)
		if err != nil {
			log.Printf("CreateMarketplaceItemMappingError: %s", err.Error())
		}
	}

	nowUnixTime := time.Now().Unix()
	createdAt := valueObject.UnixTime(nowUnixTime)
	updatedAt := valueObject.UnixTime(nowUnixTime)

	marketplaceInstalledItem := entity.NewMarketplaceInstalledItem(
		marketplaceCatalogItem.Id,
		marketplaceCatalogItem.Name,
		marketplaceCatalogItem.Type,
		installMarketplaceCatalogItem.RootDirectory,
		marketplaceCatalogItem.ServiceNames,
		[]entity.Mapping{},
		marketplaceCatalogItem.AvatarUrl,
		createdAt,
		updatedAt,
	)

	modelWithoutId := true
	marketplaceInstalledItemModel, err := dbModel.MarketplaceInstalledItem{}.ToModel(
		marketplaceInstalledItem,
		modelWithoutId,
	)
	if err != nil {
		return err
	}

	err = repo.persistentDbSvc.Handler.Create(&marketplaceInstalledItemModel).Error
	if err != nil {
		return err
	}

	return nil
}
