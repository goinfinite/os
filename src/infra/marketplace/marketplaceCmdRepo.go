package marketplaceInfra

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/google/uuid"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
	runtimeInfra "github.com/speedianet/os/src/infra/runtime"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
)

const installTempDirPath = "/app/marketplace-tmp/"

type MarketplaceCmdRepo struct {
	persistentDbSvc      *internalDbInfra.PersistentDatabaseService
	marketplaceQueryRepo *MarketplaceQueryRepo
	mappingCmdRepo       *mappingInfra.MappingCmdRepo
}

func NewMarketplaceCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceCmdRepo {
	return &MarketplaceCmdRepo{
		persistentDbSvc:      persistentDbSvc,
		marketplaceQueryRepo: NewMarketplaceQueryRepo(persistentDbSvc),
		mappingCmdRepo:       mappingInfra.NewMappingCmdRepo(persistentDbSvc),
	}
}

func (repo *MarketplaceCmdRepo) installServices(
	vhostName valueObject.Fqdn,
	services []valueObject.ServiceNameWithVersion,
) error {
	serviceQueryRepo := servicesInfra.ServicesQueryRepo{}
	serviceCmdRepo := servicesInfra.ServicesCmdRepo{}

	shouldCreatePhpVirtualHost := false
	for _, serviceWithVersion := range services {
		if serviceWithVersion.Name.String() == "php-webserver" {
			shouldCreatePhpVirtualHost = true
		}

		_, err := serviceQueryRepo.GetByName(serviceWithVersion.Name)
		if err == nil {
			continue
		}

		autoCreateMapping := false
		createServiceDto := dto.NewCreateInstallableService(
			serviceWithVersion.Name, serviceWithVersion.Version, nil, nil, autoCreateMapping,
		)

		err = serviceCmdRepo.CreateInstallable(createServiceDto)
		if err != nil {
			return errors.New("InstallRequiredServiceError: " + err.Error())
		}
	}

	if shouldCreatePhpVirtualHost {
		runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo()
		return runtimeCmdRepo.CreatePhpVirtualHost(vhostName)
	}

	return nil
}

func (repo *MarketplaceCmdRepo) parseSystemDataFields(
	installDir valueObject.UnixFilePath,
	installUrlPath valueObject.UrlPath,
	installHostname valueObject.Fqdn,
	installUuid string,
) (systemDataFields []valueObject.MarketplaceInstallableItemDataField) {
	dataMap := map[string]string{
		"installDirectory":      installDir.String(),
		"installUrlPath":        installUrlPath.String(),
		"installHostname":       installHostname.String(),
		"installUuid":           installUuid,
		"installTempDir":        installTempDirPath,
		"installRandomPassword": infraHelper.GenPass(16),
	}

	for key, value := range dataMap {
		dataFieldKey, _ := valueObject.NewDataFieldName(key)
		dataFieldValue, _ := valueObject.NewDataFieldValue(value)
		dataField, _ := valueObject.NewMarketplaceInstallableItemDataField(
			dataFieldKey, dataFieldValue,
		)
		systemDataFields = append(systemDataFields, dataField)
	}

	return systemDataFields
}

func (repo *MarketplaceCmdRepo) interpolateMissingOptionalDataFields(
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
	catalogDataFields []valueObject.MarketplaceCatalogItemDataField,
) (missingDataFields []valueObject.MarketplaceInstallableItemDataField, err error) {
	receivedDataFieldsNames := map[string]interface{}{}
	for _, receivedDataField := range receivedDataFields {
		receivedDataFieldsNames[receivedDataField.Name.String()] = nil
	}

	for _, catalogDataField := range catalogDataFields {
		if catalogDataField.IsRequired {
			continue
		}

		catalogDataFieldNameStr := catalogDataField.Name.String()
		_, alreadyFilled := receivedDataFieldsNames[catalogDataFieldNameStr]
		if alreadyFilled {
			continue
		}

		missingDataField, _ := valueObject.NewMarketplaceInstallableItemDataField(
			catalogDataField.Name, *catalogDataField.DefaultValue,
		)
		missingDataFields = append(missingDataFields, missingDataField)
	}

	return missingDataFields, nil
}

func (repo *MarketplaceCmdRepo) replaceCmdStepsPlaceholders(
	cmdSteps []valueObject.MarketplaceItemCmdStep,
	dataFields []valueObject.MarketplaceInstallableItemDataField,
) (cmdStepsWithDataFields []valueObject.MarketplaceItemCmdStep, err error) {
	dataFieldsMap := map[string]string{}
	for _, dataField := range dataFields {
		dataFieldKeyStr := dataField.Name.String()
		dataFieldsMap[dataFieldKeyStr] = dataField.Value.String()
	}

	for _, cmdStep := range cmdSteps {
		cmdStepStr := cmdStep.String()
		cmdStepDataFieldPlaceholders, _ := infraHelper.GetAllRegexGroupMatches(
			cmdStepStr, `%(.*?)%`,
		)

		for _, cmdStepDataPlaceholder := range cmdStepDataFieldPlaceholders {
			dataFieldValue := dataFieldsMap[cmdStepDataPlaceholder]
			escapedDataFieldValue := shellescape.Quote(dataFieldValue)

			cmdStepWithDataFieldStr := strings.ReplaceAll(
				cmdStepStr, "%"+cmdStepDataPlaceholder+"%", escapedDataFieldValue,
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
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
) error {
	preparedCmdSteps, err := repo.replaceCmdStepsPlaceholders(
		catalogCmdSteps, receivedDataFields,
	)
	if err != nil {
		return errors.New("ParseCmdStepWithDataFieldsError: " + err.Error())
	}

	for stepIndex, cmdStep := range preparedCmdSteps {
		_, err = infraHelper.RunCmdWithSubShell(cmdStep.String())
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

	chownRecursively := true
	chownSymlinksToo := true
	err := infraHelper.UpdatePermissionsForWebServerUse(
		installDirStr, chownRecursively, chownSymlinksToo,
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
		newMappingBase, err := valueObject.NewMappingPath(urlPath.String())
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
) (mappingIds []valueObject.MappingId, err error) {
	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(repo.persistentDbSvc)
	currentMappings, err := mappingQueryRepo.ReadByHostname(hostname)
	if err != nil {
		return mappingIds, err
	}

	currentMappingsContentHashMap := map[string]entity.Mapping{}
	for _, currentMapping := range currentMappings {
		contentHash := infraHelper.GenStrongShortHash(
			currentMapping.Hostname.String() +
				currentMapping.Path.String() +
				currentMapping.MatchPattern.String() +
				currentMapping.TargetType.String(),
		)

		currentMappingsContentHashMap[contentHash] = currentMapping
	}

	for _, mapping := range catalogMappings {
		contentHash := infraHelper.GenStrongShortHash(
			hostname.String() +
				mapping.Path.String() +
				mapping.MatchPattern.String() +
				mapping.TargetType.String(),
		)
		currentMapping, alreadyExists := currentMappingsContentHashMap[contentHash]
		if alreadyExists {
			mappingIds = append(mappingIds, currentMapping.Id)
			continue
		}

		createDto := dto.NewCreateMapping(
			hostname,
			mapping.Path,
			mapping.MatchPattern,
			mapping.TargetType,
			mapping.TargetValue,
			mapping.TargetHttpResponseCode,
		)

		mappingId, err := repo.mappingCmdRepo.Create(createDto)
		if err != nil {
			log.Printf("CreateItemMappingError: %s", err.Error())
			continue
		}

		mappingIds = append(mappingIds, mappingId)
	}

	return mappingIds, nil
}

func (repo *MarketplaceCmdRepo) persistInstalledItem(
	catalogItem entity.MarketplaceCatalogItem,
	hostname valueObject.Fqdn,
	urlPath valueObject.UrlPath,
	installDir valueObject.UnixFilePath,
	installUuid string,
	mappingsId []valueObject.MappingId,
) error {
	servicesList := []string{}
	for _, service := range catalogItem.Services {
		servicesList = append(servicesList, service.String())
	}
	servicesListStr := strings.Join(servicesList, ",")

	mappingModels := []dbModel.Mapping{}
	for _, mappingId := range mappingsId {
		mappingModel := dbModel.Mapping{ID: uint(mappingId.Get())}
		mappingModels = append(mappingModels, mappingModel)
	}

	firstCatalogItemSlug := catalogItem.Slugs[0]
	installedItemModel := dbModel.MarketplaceInstalledItem{
		Name:             catalogItem.Name.String(),
		Hostname:         hostname.String(),
		Type:             catalogItem.Type.String(),
		UrlPath:          urlPath.String(),
		InstallDirectory: installDir.String(),
		InstallUuid:      installUuid,
		Services:         servicesListStr,
		Mappings:         mappingModels,
		AvatarUrl:        catalogItem.AvatarUrl.String(),
		CatalogSlug:      firstCatalogItemSlug.String(),
	}

	return repo.persistentDbSvc.Handler.Create(&installedItemModel).Error
}

func (repo *MarketplaceCmdRepo) InstallItem(
	installDto dto.InstallMarketplaceCatalogItem,
) error {
	catalogItem, err := repo.marketplaceQueryRepo.ReadCatalogItemById(*installDto.Id)
	if err != nil {
		return errors.New("MarketplaceCatalogItemNotFound")
	}

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(repo.persistentDbSvc)
	vhost, err := vhostQueryRepo.ReadByHostname(installDto.Hostname)
	if err != nil {
		return err
	}

	err = repo.installServices(installDto.Hostname, catalogItem.Services)
	if err != nil {
		return err
	}

	installUrlPath, _ := valueObject.NewUrlPath("/")
	if installDto.UrlPath != nil {
		installUrlPath = *installDto.UrlPath
	}

	installDirStr := vhost.RootDirectory.String() + installUrlPath.GetWithoutTrailingSlash()
	installDir, err := valueObject.NewUnixFilePath(installDirStr)
	if err != nil {
		return errors.New("DefineInstallDirectoryError: " + err.Error())
	}
	installDirStr = installDir.String()

	installUuid := uuid.New().String()[:16]
	installUuidWithoutHyphens := strings.Replace(installUuid, "-", "", -1)

	systemDataFields := repo.parseSystemDataFields(
		installDir, installUrlPath, installDto.Hostname,
		installUuidWithoutHyphens,
	)
	receivedDataFields := slices.Concat(installDto.DataFields, systemDataFields)

	optionalFieldsWithDefaultValues, err := repo.interpolateMissingOptionalDataFields(
		receivedDataFields, catalogItem.DataFields,
	)
	if err != nil {
		return err
	}
	receivedDataFields = slices.Concat(receivedDataFields, optionalFieldsWithDefaultValues)

	err = infraHelper.MakeDir(installDirStr)
	if err != nil {
		return errors.New("CreateInstallDirectoryError: " + err.Error())
	}

	err = infraHelper.MakeDir(installTempDirPath)
	if err != nil {
		return errors.New("CreateTmpDirectoryError: " + err.Error())
	}

	err = repo.runCmdSteps(catalogItem.InstallCmdSteps, receivedDataFields)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd("rm", "-rf", installTempDirPath)
	if err != nil {
		return errors.New("RemoveTmpDirectoryError: " + err.Error())
	}

	err = repo.updateFilesPrivileges(installDir)
	if err != nil {
		return errors.New("UpdateFilesPrivilegesError: " + err.Error())
	}

	isRootDirectory := installDir.String() == vhost.RootDirectory.String()
	if !isRootDirectory {
		catalogItem.Mappings = repo.updateMappingsBase(
			catalogItem.Mappings, installUrlPath,
		)
	}

	mappingIds, err := repo.createMappings(installDto.Hostname, catalogItem.Mappings)
	if err != nil {
		return err
	}

	return repo.persistInstalledItem(
		catalogItem,
		installDto.Hostname,
		installUrlPath,
		installDir,
		installUuidWithoutHyphens,
		mappingIds,
	)
}

func (repo *MarketplaceCmdRepo) uninstallFilesRemoval(
	installedId valueObject.MarketplaceItemId,
) error {
	installedItem, err := repo.marketplaceQueryRepo.ReadInstalledItemById(
		installedId,
	)
	if err != nil {
		return err
	}

	trashDirName := fmt.Sprintf(
		"%s/%s-%s-%s",
		useCase.TrashDirPath,
		installedItem.AppSlug.String(),
		installedItem.Hostname.String(),
		installedItem.InstallUuid.String(),
	)
	err = infraHelper.MakeDir(trashDirName)
	if err != nil {
		return errors.New("CreateTrashDirError: " + err.Error())
	}

	catalogItem, err := repo.marketplaceQueryRepo.ReadCatalogItemBySlug(
		installedItem.AppSlug,
	)
	if err != nil {
		return err
	}

	if len(catalogItem.UninstallFilesToRemove) == 0 {
		return nil
	}

	firstFileToRemove := catalogItem.UninstallFilesToRemove[0]
	fileNameFilterParams := "-name \"" + firstFileToRemove.String() + "\""

	filesToRemoveWithoutFirstOne := catalogItem.UninstallFilesToRemove[1:]
	for _, fileToRemove := range filesToRemoveWithoutFirstOne {
		fileNameFilterParams += " -o -name \"" + fileToRemove.String() + "\""
	}

	removeFilesCmd := fmt.Sprintf(
		"find %s \\( %s \\) -maxdepth 1 -exec mv -t %s {} +",
		installedItem.InstallDirectory.String(),
		fileNameFilterParams,
		trashDirName,
	)
	_, err = infraHelper.RunCmdWithSubShell(removeFilesCmd)
	if err != nil {
		return errors.New("RemoveFilesDuringUninstallError: " + err.Error())
	}

	return nil
}

func (repo *MarketplaceCmdRepo) uninstallUnusedServices(
	servicesToUninstall []valueObject.ServiceNameWithVersion,
) error {
	serviceNamesToUninstallMap := map[string]interface{}{}
	for _, serviceNameWithVersion := range servicesToUninstall {
		serviceNamesToUninstallMap[serviceNameWithVersion.Name.String()] = nil
	}

	installedItems, err := repo.marketplaceQueryRepo.ReadInstalledItems()
	if err != nil {
		return errors.New("ReadInstalledItemsError: " + err.Error())
	}

	serviceNamesInUseMap := map[string]interface{}{}
	for _, installedItem := range installedItems {
		for _, serviceNameWithVersion := range installedItem.Services {
			serviceNamesInUseMap[serviceNameWithVersion.Name.String()] = nil
		}
	}

	unusedServiceNames := []valueObject.ServiceName{}
	for serviceNameStr := range serviceNamesToUninstallMap {
		_, isServiceInUse := serviceNamesInUseMap[serviceNameStr]
		if isServiceInUse {
			continue
		}

		serviceName, _ := valueObject.NewServiceName(serviceNameStr)
		unusedServiceNames = append(unusedServiceNames, serviceName)
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
	installedItem, err := repo.marketplaceQueryRepo.ReadInstalledItemById(
		deleteDto.InstalledId,
	)
	if err != nil {
		return err
	}

	for _, installedItemMapping := range installedItem.Mappings {
		err = repo.mappingCmdRepo.Delete(installedItemMapping.Id)
		if err != nil {
			log.Printf(
				"DeleteMappingError (%s): %s", installedItemMapping.Path, err.Error(),
			)
			continue
		}
	}

	catalogItem, err := repo.marketplaceQueryRepo.ReadCatalogItemBySlug(
		installedItem.AppSlug,
	)
	if err != nil {
		return err
	}

	err = repo.uninstallFilesRemoval(deleteDto.InstalledId)
	if err != nil {
		return err
	}

	systemDataFields := repo.parseSystemDataFields(
		installedItem.InstallDirectory, installedItem.UrlPath, installedItem.Hostname,
		installedItem.InstallUuid.String(),
	)
	err = repo.runCmdSteps(catalogItem.UninstallCmdSteps, systemDataFields)
	if err != nil {
		return err
	}

	installedItemModel := dbModel.MarketplaceInstalledItem{
		ID: uint(deleteDto.InstalledId.Get()),
	}
	err = repo.persistentDbSvc.Handler.Delete(&installedItemModel).Error
	if err != nil {
		return err
	}

	if deleteDto.ShouldUninstallServices {
		err = repo.uninstallUnusedServices(installedItem.Services)
		if err != nil {
			return err
		}
	}

	return nil
}
