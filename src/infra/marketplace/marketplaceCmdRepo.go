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
	runtimeInfra "github.com/speedianet/os/src/infra/runtime"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
)

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

func (repo *MarketplaceCmdRepo) createRequiredServices(
	vhostHostname valueObject.Fqdn,
	serviceNames []valueObject.ServiceName,
) error {
	serviceQueryRepo := servicesInfra.ServicesQueryRepo{}
	serviceCmdRepo := servicesInfra.ServicesCmdRepo{}

	shouldCreatePhpVirtualHost := false
	for _, serviceName := range serviceNames {
		if serviceName.String() == "php-webserver" {
			shouldCreatePhpVirtualHost = true
		}

		_, err := serviceQueryRepo.GetByName(serviceName)
		if err == nil {
			continue
		}

		autoCreateMapping := false
		createServiceDto := dto.NewCreateInstallableService(
			serviceName, nil, nil, nil, autoCreateMapping,
		)

		err = serviceCmdRepo.CreateInstallable(createServiceDto)
		if err != nil {
			return errors.New("InstallRequiredServiceError: " + err.Error())
		}
	}

	if shouldCreatePhpVirtualHost {
		runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo()
		err := runtimeCmdRepo.CreatePhpVirtualHost(vhostHostname)
		if err != nil {
			return err
		}
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
		"installDirectory": installDir.String(),
		"installUrlPath":   installUrlPath.String(),
		"installHostname":  installHostname.String(),
		"installUuid":      installUuid,
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
			cmdStepWithDataFieldStr := strings.ReplaceAll(
				cmdStepStr, "%"+cmdStepDataPlaceholder+"%", dataFieldValue,
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
	requiredSvcNamesListStr := []string{}
	for _, svcName := range catalogItem.RequiredServiceNames {
		requiredSvcNamesListStr = append(requiredSvcNamesListStr, svcName.String())
	}
	requiredSvcNamesStr := strings.Join(requiredSvcNamesListStr, ",")

	mappingModels := []dbModel.Mapping{}
	for _, mappingId := range mappingsId {
		mappingModel := dbModel.Mapping{ID: uint(mappingId.Get())}
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

	return repo.persistentDbSvc.Handler.Create(&installedItemModel).Error
}

func (repo *MarketplaceCmdRepo) InstallItem(
	installDto dto.InstallMarketplaceCatalogItem,
) error {
	catalogItem, err := repo.marketplaceQueryRepo.ReadCatalogItemById(installDto.Id)
	if err != nil {
		return errors.New("MarketplaceCatalogItemNotFound")
	}

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhost, err := vhostQueryRepo.GetByHostname(installDto.Hostname)
	if err != nil {
		return err
	}

	err = repo.createRequiredServices(installDto.Hostname, catalogItem.RequiredServiceNames)
	if err != nil {
		return err
	}

	installUrlPath, _ := valueObject.NewUrlPath("/")
	if installDto.UrlPath != nil {
		installUrlPath = *installDto.UrlPath
	}

	installDirStr := vhost.RootDirectory.String() + installUrlPath.GetWithoutLeadingSlash()
	installDir, err := valueObject.NewUnixFilePath(installDirStr)
	if err != nil {
		return errors.New("DefineInstallDirectoryError: " + err.Error())
	}

	installUuid := uuid.New().String()[:16]
	installUuidWithoutHyphens := strings.Replace(installUuid, "-", "", -1)

	systemDataFields := repo.parseSystemDataFields(
		installDir, installUrlPath, installDto.Hostname, installUuidWithoutHyphens,
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

	err = repo.runCmdSteps(catalogItem.CmdSteps, receivedDataFields)
	if err != nil {
		return err
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

func (repo *MarketplaceCmdRepo) getServiceNamesInUse() (
	[]valueObject.ServiceName, error,
) {
	servicesInUse := []valueObject.ServiceName{}

	installedItems, err := repo.marketplaceQueryRepo.ReadInstalledItems()
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
