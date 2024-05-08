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

func (repo *MarketplaceCmdRepo) parseSystemDataFields(
	installDir valueObject.UnixFilePath,
	installUrlPath valueObject.UrlPath,
	installHostname valueObject.Fqdn,
	installUuid string,
) []valueObject.MarketplaceInstallableItemDataField {
	systemDataFields := []valueObject.MarketplaceInstallableItemDataField{}

	installDirDataFieldKey, _ := valueObject.NewDataFieldName("installDirectory")
	installDirDataFieldValue, _ := valueObject.NewDataFieldValue(installDir.String())
	installDirDataField, _ := valueObject.NewMarketplaceInstallableItemDataField(
		installDirDataFieldKey,
		installDirDataFieldValue,
	)

	installUrlPathDataFieldKey, _ := valueObject.NewDataFieldName("installUrlPath")
	installUrlPathDataFieldValue, _ := valueObject.NewDataFieldValue(installUrlPath.String())
	installUrlPathDataField, _ := valueObject.NewMarketplaceInstallableItemDataField(
		installUrlPathDataFieldKey,
		installUrlPathDataFieldValue,
	)

	installHostnameDataFieldKey, _ := valueObject.NewDataFieldName("installHostname")
	installHostnameDataFieldValue, _ := valueObject.NewDataFieldValue(
		installHostname.String(),
	)
	installHostnameDataField, _ := valueObject.NewMarketplaceInstallableItemDataField(
		installHostnameDataFieldKey,
		installHostnameDataFieldValue,
	)

	installUuidDataFieldKey, _ := valueObject.NewDataFieldName("installUuid")
	installUuidDataFieldValue, _ := valueObject.NewDataFieldValue(installUuid)
	installUuidDataField, _ := valueObject.NewMarketplaceInstallableItemDataField(
		installUuidDataFieldKey,
		installUuidDataFieldValue,
	)

	return append(
		systemDataFields,
		installDirDataField,
		installUrlPathDataField,
		installHostnameDataField,
		installUuidDataField,
	)
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
			receivedDataField.Name.String(),
		)
	}

	for _, catalogDataField := range catalogDataFields {
		if catalogDataField.IsRequired {
			continue
		}

		catalogDataFieldKeyStr := catalogDataField.Name.String()
		catalogDataFieldAlreadyFilled := slices.Contains(
			receivedDataFieldsKeys,
			catalogDataFieldKeyStr,
		)
		if catalogDataFieldAlreadyFilled {
			continue
		}

		catalogDataFieldAsInstallable, _ := valueObject.NewMarketplaceInstallableItemDataField(
			catalogDataField.Name,
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
		dataFieldKeyStr := dataField.Name.String()
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

func (repo *MarketplaceCmdRepo) updateFilesPrivileges(
	installDir valueObject.UnixFilePath,
) error {
	installDirStr := installDir.String()
	_, err := infraHelper.RunCmdWithSubShell(
		"chown -R nobody:nogroup -L " + installDirStr,
	)
	if err != nil {
		return errors.New("ChownError (" + installDirStr + "): " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell(
		`find ` + installDirStr + ` -type d -exec chmod 755 {} \; && find ` +
			installDirStr + ` -type f -exec chmod 644 {} \;`,
	)
	if err != nil {
		return errors.New("ChmodError (" + installDirStr + "): " + err.Error())
	}

	return nil
}

func (repo *MarketplaceCmdRepo) updateMappingsBase(
	catalogMappings []valueObject.MarketplaceItemMapping,
	urlPath valueObject.UrlPath,
) []valueObject.MarketplaceItemMapping {
	for mappingIndex, catalogMapping := range catalogMappings {
		isPathRoot := catalogMapping.Path.String() == "/"
		if !isPathRoot {
			continue
		}

		catalogMappingWithNewPath := catalogMapping
		newMappingBase, err := valueObject.NewMappingPath(
			urlPath.String(),
		)
		if err != nil {
			log.Printf("%s: %s", err.Error(), urlPath.String())
			continue
		}
		catalogMappingWithNewPath.Path = newMappingBase

		catalogMappings[mappingIndex] = catalogMappingWithNewPath
	}

	return catalogMappings
}

func (repo *MarketplaceCmdRepo) createMappings(
	hostname valueObject.Fqdn,
	catalogMappings []valueObject.MarketplaceItemMapping,
) (createdMappings []entity.Mapping, err error) {
	for _, catalogMapping := range catalogMappings {
		createCatalogItemMapping := dto.NewCreateMapping(
			hostname,
			catalogMapping.Path,
			catalogMapping.MatchPattern,
			catalogMapping.TargetType,
			catalogMapping.TargetValue,
			catalogMapping.TargetHttpResponseCode,
		)

		mappingCmdRepo := mappingInfra.NewMappingCmdRepo(repo.persistentDbSvc)
		mappingId, err := mappingCmdRepo.Create(createCatalogItemMapping)
		if err != nil {
			log.Printf("CreateMarketplaceItemMappingError: %s", err.Error())
		}

		createdMapping := entity.NewMapping(
			mappingId, hostname, catalogMapping.Path, catalogMapping.MatchPattern,
			catalogMapping.TargetType, catalogMapping.TargetValue,
			catalogMapping.TargetHttpResponseCode,
		)
		createdMappings = append(createdMappings, createdMapping)
	}

	return createdMappings, nil
}

func (repo *MarketplaceCmdRepo) persistInstalledItem(
	catalogItem entity.MarketplaceCatalogItem,
	hostname valueObject.Fqdn,
	urlPath valueObject.UrlPath,
	installDir valueObject.UnixFilePath,
	installUuid string,
	createdMappings []entity.Mapping,
) error {
	requiredSvcNamesListStr := []string{}
	for _, svcName := range catalogItem.RequiredServiceNames {
		requiredSvcNamesListStr = append(requiredSvcNamesListStr, svcName.String())
	}
	requiredSvcNamesStr := strings.Join(requiredSvcNamesListStr, ",")

	mappingModels := []dbModel.Mapping{}
	for _, createdMapping := range createdMappings {
		mappingModel := dbModel.Mapping{}.ToModel(createdMapping)
		mappingModels = append(mappingModels, mappingModel)
	}

	installedItemModel := dbModel.MarketplaceInstalledItem{
		Name:                 catalogItem.Name.String(),
		Hostname:             hostname.String(),
		Type:                 catalogItem.Type.String(),
		UrlPath:              urlPath.String(),
		InstallDirectory:     installDir.String(),
		InstallUuid:          installUuid,
		RequiredServiceNames: requiredSvcNamesStr,
		Mappings:             mappingModels,
		AvatarUrl:            catalogItem.AvatarUrl.String(),
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

	err = repo.createRequiredServices(catalogItem.RequiredServiceNames)
	if err != nil {
		return err
	}

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhost, err := vhostQueryRepo.GetByHostname(installDto.Hostname)
	if err != nil {
		return err
	}

	installUrlPath, _ := valueObject.NewUrlPath("/")
	if installDto.UrlPath != nil {
		installUrlPath = *installDto.UrlPath
	}

	installDirStr := vhost.RootDirectory.String() + installUrlPath.GetWithoutLeadingSlash()
	installDir, _ := valueObject.NewUnixFilePath(installDirStr)

	installUuid := uuid.New().String()[:16]
	installUuidWithoutHyphens := strings.Replace(installUuid, "-", "", -1)

	systemDataFields := repo.parseSystemDataFields(
		installDir,
		installUrlPath,
		installDto.Hostname,
		installUuidWithoutHyphens,
	)
	receivedDataFields := slices.Concat(installDto.DataFields, systemDataFields)

	err = repo.runCmdSteps(
		catalogItem.CmdSteps,
		catalogItem.DataFields,
		receivedDataFields,
	)
	if err != nil {
		return err
	}

	err = repo.updateFilesPrivileges(installDir)
	if err != nil {
		return errors.New("UpdateFilesPrivileges: " + err.Error())
	}

	isRootDirectory := installDir.String() == vhost.RootDirectory.String()
	if !isRootDirectory {
		catalogItem.Mappings = repo.updateMappingsBase(
			catalogItem.Mappings,
			installUrlPath,
		)
	}

	createdMappings, err := repo.createMappings(
		installDto.Hostname, catalogItem.Mappings,
	)
	if err != nil {
		return err
	}

	return repo.persistInstalledItem(
		catalogItem,
		installDto.Hostname,
		installUrlPath,
		installDir,
		installUuidWithoutHyphens,
		createdMappings,
	)
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
			installedItem.RequiredServiceNames,
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
	deleteDto dto.DeleteMarketplaceInstalledItem,
) error {
	installedItem, err := repo.marketplaceQueryRepo.GetInstalledItemById(
		deleteDto.InstalledId,
	)
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
		ID: uint(deleteDto.InstalledId.Get()),
	}
	err = repo.persistentDbSvc.Handler.Delete(&installedItemModel).Error
	if err != nil {
		return err
	}

	if deleteDto.ShouldUninstallServices {
		err = repo.uninstallServices(installedItem.RequiredServiceNames)
		if err != nil {
			return err
		}
	}

	if deleteDto.ShouldRemoveFiles {
		installDirStr := installedItem.InstallDirectory.String()
		err = os.RemoveAll(installDirStr)
		if err != nil {
			return errors.New("DeleteInstalledItemFilesError: " + err.Error())
		}

		err = infraHelper.MakeDir(installDirStr)
		if err != nil {
			return errors.New("CreateEmptyInstallDirectoryError: " + err.Error())
		}
	}

	return nil
}
