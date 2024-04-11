package marketplaceInfra

import (
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
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

func (repo *MarketplaceCmdRepo) createRequiredServices(
	catalogRequiredSvcNames []valueObject.ServiceName,
) error {
	svcQueryRepo := servicesInfra.ServicesQueryRepo{}
	svcCmdRepo := servicesInfra.ServicesCmdRepo{}
	for _, requiredSvcName := range catalogRequiredSvcNames {
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

	return nil
}

func (repo *MarketplaceCmdRepo) addMissingOptionalDataFieldsToMap(
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

func (repo *MarketplaceCmdRepo) parseCmdStepWithReceivedDataFields(
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

func (repo *MarketplaceCmdRepo) runCmdStepsWithDataFields(
	catalogCmdSteps []valueObject.MarketplaceItemInstallStep,
	catalogDataFields []valueObject.MarketplaceCatalogItemDataField,
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
	installDir valueObject.UnixFilePath,
) error {
	receivedDataFieldsMap := map[string]string{}

	for _, receivedDataField := range receivedDataFields {
		receivedDataFieldKeyStr := receivedDataField.Key.String()
		receivedDataFieldsMap[receivedDataFieldKeyStr] = receivedDataField.Value.String()
	}

	receivedDataFieldsMap["installDirectory"] = installDir.String()
	receivedDataFieldsMap["installUuid"] = uuid.New().String()[:16]
	repo.addMissingOptionalDataFieldsToMap(
		&receivedDataFieldsMap,
		catalogDataFields,
	)

	for _, cmdStep := range catalogCmdSteps {
		cmdStepRequiredDataFields, err := repo.parseCmdStepWithReceivedDataFields(
			cmdStep,
			receivedDataFieldsMap,
		)
		if err != nil {
			return errors.New("ParseCmdStepWithDataFieldsError: " + err.Error())
		}

		_, err = infraHelper.RunCmdWithSubShell(cmdStepRequiredDataFields)
		if err != nil {
			return errors.New(
				"RunCmdStepError (" + cmdStepRequiredDataFields + "): " + err.Error(),
			)
		}
	}

	return nil
}

func (repo *MarketplaceCmdRepo) createMappings(
	hostname valueObject.Fqdn,
	catalogMappings []valueObject.MarketplaceItemMapping,
) error {
	for _, catalogMapping := range catalogMappings {
		createCatalogItemMapping := dto.NewCreateMapping(
			hostname,
			catalogMapping.Path,
			catalogMapping.MatchPattern,
			catalogMapping.TargetType,
			catalogMapping.TargetServiceName,
			catalogMapping.TargetUrl,
			catalogMapping.TargetHttpResponseCode,
			catalogMapping.TargetInlineHtmlContent,
		)

		vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}
		err := vhostCmdRepo.CreateMapping(createCatalogItemMapping)
		if err != nil {
			log.Printf("CreateMarketplaceItemMappingError: %s", err.Error())
		}
	}

	return nil
}

func (repo *MarketplaceCmdRepo) persistInstalledItem(
	catalogItem entity.MarketplaceCatalogItem,
	installDir valueObject.UnixFilePath,
) error {
	svcNamesListStr := []string{}
	for _, svcName := range catalogItem.ServiceNames {
		svcNamesListStr = append(svcNamesListStr, svcName.String())
	}
	svcNamesStr := strings.Join(svcNamesListStr, ",")

	installedItemModel := dbModel.MarketplaceInstalledItem{
		Name:             catalogItem.Name.String(),
		Type:             catalogItem.Type.String(),
		InstallDirectory: installDir.String(),
		ServiceNames:     svcNamesStr,
		AvatarUrl:        catalogItem.AvatarUrl.String(),
	}

	err := repo.persistentDbSvc.Handler.Create(&installedItemModel).Error
	if err != nil {
		return err
	}

	return nil
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

	err = repo.createRequiredServices(catalogItem.ServiceNames)
	if err != nil {
		return err
	}

	installDirStr := infraData.GlobalConfigs.PrimaryPublicDir
	if installDto.InstallDirectory != nil {
		vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
		vhost, err := vhostQueryRepo.GetByHostname(installDto.Hostname)
		if err != nil {
			return err
		}

		installDirStr = installDto.InstallDirectory.String()
		hasLeadingSlash := strings.HasPrefix(installDirStr, "/")
		if !hasLeadingSlash {
			installDirStr = "/" + installDirStr
		}

		installDirStr = vhost.RootDirectory.String() + installDirStr
	}
	installDir, _ := valueObject.NewUnixFilePath(installDirStr)

	err = repo.runCmdStepsWithDataFields(
		catalogItem.CmdSteps,
		catalogItem.DataFields,
		installDto.DataFields,
		installDir,
	)
	if err != nil {
		return err
	}

	err = repo.createMappings(installDto.Hostname, catalogItem.Mappings)
	if err != nil {
		return err
	}

	return repo.persistInstalledItem(catalogItem, installDir)
}
