package marketplaceInfra

import (
	"errors"
	"log"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	"github.com/speedianet/os/src/infra/infraData"
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

func (repo *MarketplaceCmdRepo) getReceivedDataFieldsAsMap(
	receivedDataFields []valueObject.MarketplaceInstalledItemDataField,
) map[string]string {
	receivedDataFieldsMap := map[string]string{}

	for _, receivedDataField := range receivedDataFields {
		receivedDataFieldKeyStr := receivedDataField.Key.String()
		receivedDataFieldsMap[receivedDataFieldKeyStr] = receivedDataField.Value.String()
	}

	return receivedDataFieldsMap
}

func (repo *MarketplaceCmdRepo) addMissingOptionalDataFields(
	receivedDataFieldsMap *map[string]string,
	requiredDataFields []valueObject.MarketplaceCatalogItemDataField,
) {
	for _, requiredDataField := range requiredDataFields {
		if requiredDataField.IsRequired {
			continue
		}

		requiredKeyStr := requiredDataField.Key.String()
		if len((*receivedDataFieldsMap)[requiredKeyStr]) != 0 {
			continue
		}

		requiredDefaultValueStr := requiredDataField.DefaultValue.String()
		(*receivedDataFieldsMap)[requiredKeyStr] = requiredDefaultValueStr
	}
}

func (repo *MarketplaceCmdRepo) getCmdStepWithReceivedDataFields(
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

	installDirStr := infraData.GlobalConfigs.PrimaryPublicDir
	if installDto.InstallDirectory != nil {
		vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
		vhost, err := vhostQueryRepo.GetByHostname(installDto.Hostname)
		if err != nil {
			return err
		}

		installDirStr = vhost.RootDirectory.String() + installDto.InstallDirectory.String()
	}

	receivedDataFielsdMap := repo.getReceivedDataFieldsAsMap(installDto.DataFields)
	receivedDataFielsdMap["installDirectory"] = installDirStr
	repo.addMissingOptionalDataFields(&receivedDataFielsdMap, catalogItem.DataFields)
	log.Printf("receivedDataFielsdMap: %+v", receivedDataFielsdMap)

	for _, cmdStep := range catalogItem.CmdSteps {
		cmdStepRequiredDataFields, err := repo.getCmdStepWithReceivedDataFields(
			cmdStep,
			receivedDataFielsdMap,
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

	installDir, _ := valueObject.NewUnixFilePath(installDirStr)
	installedItemDto := dto.CreateNewMarketplaceInstalledItem(
		catalogItem.Name,
		catalogItem.Type,
		installDir,
		catalogItem.ServiceNames,
		[]entity.Mapping{},
		catalogItem.AvatarUrl,
	)

	installedItemModel, err := dbModel.MarketplaceInstalledItem{}.ToModelFromDto(
		installedItemDto,
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
