package marketplaceInfra

import (
	"errors"
	"log"
	"os"
	"slices"
	"strconv"
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
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
)

type MarketplaceCmdRepo struct {
	persistentDbSvc      *internalDbInfra.PersistentDatabaseService
	marketplaceQueryRepo *MarketplaceQueryRepo
}

func NewMarketplaceCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceCmdRepo {
	marketplaceQueryRepo := NewMarketplaceQueryRepo(persistentDbSvc)

	return &MarketplaceCmdRepo{
		persistentDbSvc:      persistentDbSvc,
		marketplaceQueryRepo: marketplaceQueryRepo,
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

func (repo *MarketplaceCmdRepo) interpolateMissingDataFields(
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
	catalogDataFields []valueObject.MarketplaceCatalogItemDataField,
) ([]valueObject.MarketplaceInstallableItemDataField, error) {
	missingCatalogOptionalDataFields := []valueObject.MarketplaceInstallableItemDataField{}

	receivedDataFieldsKeys := []string{}
	for _, receivedDataField := range receivedDataFields {
		receivedDataFieldsKeys = append(
			receivedDataFieldsKeys,
			receivedDataField.Key.String(),
		)
	}

	for _, catalogDataField := range catalogDataFields {
		if catalogDataField.IsRequired {
			continue
		}

		catalogDataFieldKeyStr := catalogDataField.Key.String()
		catalogDataFieldAlreadyFilled := slices.Contains(
			receivedDataFieldsKeys,
			catalogDataFieldKeyStr,
		)
		if catalogDataFieldAlreadyFilled {
			continue
		}

		catalogDataFieldAsInstallable, _ := valueObject.NewMarketplaceInstallableItemDataField(
			catalogDataField.Key,
			*catalogDataField.DefaultValue,
		)
		missingCatalogOptionalDataFields = append(
			missingCatalogOptionalDataFields,
			catalogDataFieldAsInstallable,
		)
	}

	return missingCatalogOptionalDataFields, nil
}

func (repo *MarketplaceCmdRepo) replaceCmdStepsPlaceholders(
	cmdSteps []valueObject.MarketplaceItemCmdStep,
	dataFields []valueObject.MarketplaceInstallableItemDataField,
) ([]valueObject.MarketplaceItemCmdStep, error) {
	cmdStepsWithDataFields := []valueObject.MarketplaceItemCmdStep{}

	dataFieldsMap := map[string]string{}
	for _, dataField := range dataFields {
		dataFieldKeyStr := dataField.Key.String()
		dataFieldsMap[dataFieldKeyStr] = dataField.Value.String()
	}

	for _, cmdStep := range cmdSteps {
		cmdStepStr := cmdStep.String()
		cmdStepDataFieldKeys, _ := infraHelper.GetAllRegexGroupMatches(
			cmdStepStr,
			`%(.*?)%`,
		)

		for _, cmdStepDataFieldKey := range cmdStepDataFieldKeys {
			dataFieldValue := dataFieldsMap[cmdStepDataFieldKey]
			cmdStepWithDataFieldStr := strings.ReplaceAll(
				cmdStepStr,
				"%"+cmdStepDataFieldKey+"%",
				dataFieldValue,
			)
			cmdStepStr = cmdStepWithDataFieldStr
		}

		cmdStepWithDataField, _ := valueObject.NewMarketplaceItemCmdStep(cmdStepStr)
		cmdStepsWithDataFields = append(cmdStepsWithDataFields, cmdStepWithDataField)
	}

	return cmdStepsWithDataFields, nil
}

func (repo *MarketplaceCmdRepo) runCmdSteps(
	catalogCmdSteps []valueObject.MarketplaceItemCmdStep,
	catalogDataFields []valueObject.MarketplaceCatalogItemDataField,
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
) error {
	missingCatalogOptionalDataFields, err := repo.interpolateMissingDataFields(
		receivedDataFields,
		catalogDataFields,
	)
	if err != nil {
		return err
	}
	receivedDataFields = slices.Concat(
		receivedDataFields,
		missingCatalogOptionalDataFields,
	)

	preparedCmdSteps, err := repo.replaceCmdStepsPlaceholders(
		catalogCmdSteps,
		receivedDataFields,
	)
	if err != nil {
		return errors.New("ParseCmdStepWithDataFieldsError: " + err.Error())
	}

	for stepIndex, cmdStep := range preparedCmdSteps {
		cmdStepStr := cmdStep.String()
		_, err = infraHelper.RunCmdWithSubShell(cmdStepStr)
		if err != nil {
			stepIndexStr := strconv.Itoa(stepIndex)
			return errors.New(
				"RunCmdStepError (" + stepIndexStr + "): " + err.Error(),
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

		mappingCmdRepo := mappingInfra.NewMappingCmdRepo(repo.persistentDbSvc)
		_, err := mappingCmdRepo.Create(createCatalogItemMapping)
		if err != nil {
			log.Printf("CreateMarketplaceItemMappingError: %s", err.Error())
		}
	}

	return nil
}

func (repo *MarketplaceCmdRepo) persistInstalledItem(
	catalogItem entity.MarketplaceCatalogItem,
	installDir valueObject.UnixFilePath,
	installUuid string,
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
		InstallUuid:      installUuid,
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
	catalogItem, err := repo.marketplaceQueryRepo.GetCatalogItemById(
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

	installDirDataFieldKey, _ := valueObject.NewDataFieldKey("installDirectory")
	installDirDataFieldValue, _ := valueObject.NewDataFieldValue(installDir.String())
	installDirDataField, _ := valueObject.NewMarketplaceInstallableItemDataField(
		installDirDataFieldKey,
		installDirDataFieldValue,
	)

	installUuid := uuid.New().String()[:16]
	installUuidDataFieldKey, _ := valueObject.NewDataFieldKey("installUuid")
	installUuidDataFieldValue, _ := valueObject.NewDataFieldValue(installUuid)
	installUuidDataField, _ := valueObject.NewMarketplaceInstallableItemDataField(
		installUuidDataFieldKey,
		installUuidDataFieldValue,
	)

	receivedDataFields := installDto.DataFields
	receivedDataFields = append(
		receivedDataFields,
		installDirDataField,
		installUuidDataField,
	)

	err = repo.runCmdSteps(
		catalogItem.CmdSteps,
		catalogItem.DataFields,
		receivedDataFields,
	)
	if err != nil {
		return err
	}

	err = repo.createMappings(installDto.Hostname, catalogItem.Mappings)
	if err != nil {
		return err
	}

	return repo.persistInstalledItem(catalogItem, installDir, installUuid)
}

func (repo *MarketplaceCmdRepo) getServiceNamesInUse() (
	[]valueObject.ServiceName, error,
) {
	servicesInUse := []valueObject.ServiceName{}

	installedItems, err := repo.marketplaceQueryRepo.GetInstalledItems()
	if err != nil {
		return servicesInUse, err
	}

	for _, installedItem := range installedItems {
		servicesInUse = slices.Concat(
			servicesInUse,
			installedItem.ServiceNames,
		)
	}

	return servicesInUse, nil
}

func (repo *MarketplaceCmdRepo) uninstallServices(
	installedServiceNames []valueObject.ServiceName,
) error {
	serviceNamesInUse, err := repo.getServiceNamesInUse()
	if err != nil {
		return err
	}

	unusedServiceNames := []valueObject.ServiceName{}
	for _, installedServiceName := range installedServiceNames {
		isInstalledServiceInUse := slices.Contains(
			serviceNamesInUse, installedServiceName,
		)
		if isInstalledServiceInUse {
			continue
		}

		unusedServiceNames = append(unusedServiceNames, installedServiceName)
	}

	servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
	for _, unusedService := range unusedServiceNames {
		err = servicesCmdRepo.Uninstall(unusedService)
		if err != nil {
			log.Printf("UninstallUnusedServiceError: %s", err.Error())
			continue
		}
	}

	return nil
}

func (repo *MarketplaceCmdRepo) UninstallItem(
	installedId valueObject.MarketplaceInstalledItemId,
	shouldUninstallServices bool,
) error {
	installedItem, err := repo.marketplaceQueryRepo.GetInstalledItemById(installedId)
	if err != nil {
		return err
	}

	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(repo.persistentDbSvc)
	for _, installedItemMapping := range installedItem.Mappings {
		err = mappingCmdRepo.Delete(installedItemMapping.Id)
		if err != nil {
			log.Printf(
				"DeleteInstalledItemMappingError (%s): %s",
				installedItemMapping.Path,
				err.Error(),
			)
			continue
		}
	}

	installedItemModel := dbModel.MarketplaceInstalledItem{
		ID: uint(installedId.Get()),
	}
	err = repo.persistentDbSvc.Handler.Delete(&installedItemModel).Error
	if err != nil {
		return err
	}

	if shouldUninstallServices {
		err = repo.uninstallServices(installedItem.ServiceNames)
		if err != nil {
			return err
		}
	}

	installDirStr := installedItem.InstallDirectory.String()
	err = os.RemoveAll(installDirStr)
	if err != nil {
		return errors.New("DeleteInstalledItemFilesError: " + err.Error())
	}

	err = infraHelper.MakeDir(installDirStr)
	if err != nil {
		return errors.New("CreateEmptyInstallDirectoryError: " + err.Error())
	}

	return nil
}
