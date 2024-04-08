package marketplaceInfra

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
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
	cmdStepRequiredDataFields, _ := infraHelper.GetAllRegexGroupMatches(
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

func (repo *MarketplaceCmdRepo) moveInstalledItem(
	installedItemName valueObject.MarketplaceItemName,
	rootDirectory valueObject.UnixFilePath,
) error {
	installedItemSrcPath := "/speedia/" + installedItemName.String()
	_, err := infraHelper.RunCmd(
		"mv",
		installedItemSrcPath,
		rootDirectory.String(),
	)

	return err
}

func (repo *MarketplaceCmdRepo) InstallItem(
	installDto dto.InstallMarketplaceCatalogItem,
) error {
	catalogItem, err := repo.queryRepo.GetCatalogItemById(
		installDto.Id,
	)
	if err != nil {
		return errors.New("MarketplaceCatalogItemNotFound")
	}

	svcQueryRepo := servicesInfra.ServicesQueryRepo{}
	svcCmdRepo := servicesInfra.ServicesCmdRepo{}
	for _, requiredSvcName := range catalogItem.ServiceNames {
		_, err := svcQueryRepo.GetByName(requiredSvcName)
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

		err = svcCmdRepo.CreateInstallable(requiredService)
		if err != nil {
			return errors.New("InstallRequiredServiceError: " + err.Error())
		}
	}

	dataFieldsMap := repo.getDataFieldsAsMap(installDto.DataFields)
	for _, cmdStep := range catalogItem.CmdSteps {
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

	err = repo.moveInstalledItem(
		catalogItem.Name,
		installDto.RootDirectory,
	)
	if err != nil {
		return err
	}

	for _, catalogItemMapping := range catalogItem.Mappings {
		createCatalogItemMapping := dto.NewCreateMapping(
			installDto.Hostname,
			catalogItemMapping.Path,
			catalogItemMapping.MatchPattern,
			catalogItemMapping.TargetType,
			catalogItemMapping.TargetServiceName,
			catalogItemMapping.TargetUrl,
			catalogItemMapping.TargetHttpResponseCode,
			catalogItemMapping.TargetInlineHtmlContent,
		)

		vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}
		err = vhostCmdRepo.CreateMapping(createCatalogItemMapping)
		if err != nil {
			log.Printf("CreateMarketplaceItemMappingError: %s", err.Error())
		}
	}

	nowUnixTime := time.Now().Unix()
	createdAt := valueObject.UnixTime(nowUnixTime)
	updatedAt := valueObject.UnixTime(nowUnixTime)

	installedItemEntity := entity.NewMarketplaceInstalledItem(
		catalogItem.Id,
		catalogItem.Name,
		catalogItem.Type,
		installDto.RootDirectory,
		catalogItem.ServiceNames,
		[]entity.Mapping{},
		catalogItem.AvatarUrl,
		createdAt,
		updatedAt,
	)

	modelWithoutId := true
	installedItemModel, err := dbModel.MarketplaceInstalledItem{}.ToModel(
		installedItemEntity,
		modelWithoutId,
	)
	if err != nil {
		return err
	}

	err = repo.persistentDbSvc.Handler.Create(&installedItemModel).Error
	if err != nil {
		return err
	}

	return nil
}
