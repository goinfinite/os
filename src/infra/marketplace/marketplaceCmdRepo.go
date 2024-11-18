package marketplaceInfra

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	runtimeInfra "github.com/goinfinite/os/src/infra/runtime"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	mappingInfra "github.com/goinfinite/os/src/infra/vhost/mapping"
	"github.com/google/uuid"
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
	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(repo.persistentDbSvc)
	serviceCmdRepo := servicesInfra.NewServicesCmdRepo(repo.persistentDbSvc)

	shouldCreatePhpVirtualHost := false
	for _, serviceWithVersion := range services {
		if serviceWithVersion.Name.String() == "php-webserver" {
			shouldCreatePhpVirtualHost = true
		}

		shouldIncludeMetrics := false
		readInstalledServiceDto := dto.ReadInstalledServicesItemsRequest{
			ServiceName:          &serviceWithVersion.Name,
			ShouldIncludeMetrics: &shouldIncludeMetrics,
		}
		_, err := servicesQueryRepo.ReadOneInstalledItem(readInstalledServiceDto)
		if err != nil {
			return err
		}

		createServiceDto := dto.NewCreateInstallableService(
			serviceWithVersion.Name, []valueObject.ServiceEnv{},
			[]valueObject.PortBinding{}, serviceWithVersion.Version,
			nil, nil, nil, nil, nil, nil,
		)

		_, err = serviceCmdRepo.CreateInstallable(createServiceDto)
		if err != nil {
			return errors.New("InstallRequiredServiceError: " + err.Error())
		}
	}

	if shouldCreatePhpVirtualHost {
		runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo(repo.persistentDbSvc)
		return runtimeCmdRepo.CreatePhpVirtualHost(vhostName)
	}

	return nil
}

func (repo *MarketplaceCmdRepo) parseSystemDataFields(
	installDir valueObject.UnixFilePath,
	installUrlPath valueObject.UrlPath,
	installHostname valueObject.Fqdn,
	installUuid valueObject.MarketplaceInstalledItemUuid,
) (systemDataFields []valueObject.MarketplaceInstallableItemDataField) {
	dataMap := map[string]string{
		"installDirectory":      installDir.String(),
		"installUrlPath":        installUrlPath.String(),
		"installHostname":       installHostname.String(),
		"installUuid":           installUuid.String(),
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
	cmdSteps []valueObject.UnixCommand,
	dataFields []valueObject.MarketplaceInstallableItemDataField,
) (cmdStepsWithDataFields []valueObject.UnixCommand, err error) {
	dataFieldsMap := map[string]string{}
	for _, dataField := range dataFields {
		dataFieldKeyStr := dataField.Name.String()
		dataFieldsMap[dataFieldKeyStr] = dataField.Value.String()
	}

	for _, cmdStep := range cmdSteps {
		cmdStepStr := cmdStep.String()
		cmdStepDataFieldPlaceholders := infraHelper.GetAllRegexGroupMatches(
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

		cmdStepWithDataField, _ := valueObject.NewUnixCommand(cmdStepStr)
		cmdStepsWithDataFields = append(cmdStepsWithDataFields, cmdStepWithDataField)
	}

	return cmdStepsWithDataFields, nil
}

func (repo *MarketplaceCmdRepo) runCmdSteps(
	stepType string,
	catalogCmdSteps []valueObject.UnixCommand,
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
) error {
	preparedCmdSteps, err := repo.replaceCmdStepsPlaceholders(
		catalogCmdSteps, receivedDataFields,
	)
	if err != nil {
		return errors.New("ParseCmdStepWithDataFieldsError: " + err.Error())
	}

	for stepIndex, cmdStep := range preparedCmdSteps {
		stepStr := cmdStep.String()

		slog.Debug("Running"+stepType+"Step", slog.String("step", stepStr))

		stepOutput, err := infraHelper.RunCmdWithSubShell(stepStr)
		if err != nil {
			stepIndexStr := strconv.Itoa(stepIndex)
			combinedOutput := stepOutput + " " + err.Error()
			return errors.New(
				stepType + "CmdStepError (" + stepIndexStr + "): " + combinedOutput,
			)
		}
	}

	return nil
}

func (repo *MarketplaceCmdRepo) updateFilesPrivileges(
	targetDir valueObject.UnixFilePath,
) error {
	targetDirStr := targetDir.String()

	chownRecursively := true
	chownSymlinksToo := true
	err := infraHelper.UpdatePermissionsForWebServerUse(
		targetDirStr, chownRecursively, chownSymlinksToo,
	)
	if err != nil {
		return errors.New("ChownError (" + targetDirStr + "): " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell(
		`find ` + targetDirStr + ` -type d -exec chmod 755 {} \; && find ` +
			targetDirStr + ` -type f -exec chmod 644 {} \;`,
	)
	if err != nil {
		return errors.New("ChmodError (" + targetDirStr + "): " + err.Error())
	}

	return nil
}

func (repo *MarketplaceCmdRepo) updateMappingsBase(
	catalogMappings []valueObject.MarketplaceItemMapping,
	installUrlPath valueObject.UrlPath,
) []valueObject.MarketplaceItemMapping {
	installUrlPathStr := installUrlPath.String()

	for mappingIndex, catalogMapping := range catalogMappings {
		pathStr := catalogMapping.Path.String()
		if installUrlPathStr == pathStr {
			continue
		}
		if pathStr == "/" {
			pathStr = ""
		}

		rawUpdatedPath := installUrlPathStr + pathStr
		updatedPath, err := valueObject.NewMappingPath(rawUpdatedPath)
		if err != nil {
			slog.Error(
				err.Error(),
				slog.Int("index", mappingIndex),
				slog.String("path", rawUpdatedPath),
			)
			continue
		}

		catalogMappings[mappingIndex].Path = updatedPath
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
			hostname.String() + mapping.Path.String() + mapping.MatchPattern.String() +
				mapping.TargetType.String(),
		)
		currentMapping, alreadyExists := currentMappingsContentHashMap[contentHash]
		if alreadyExists {
			mappingIds = append(mappingIds, currentMapping.Id)
			continue
		}

		createDto := dto.NewCreateMapping(
			hostname, mapping.Path, mapping.MatchPattern, mapping.TargetType,
			mapping.TargetValue, mapping.TargetHttpResponseCode,
		)

		mappingId, err := repo.mappingCmdRepo.Create(createDto)
		if err != nil {
			slog.Error("CreateItemMappingError", slog.Any("error", err))
			continue
		}

		mappingIds = append(mappingIds, mappingId)
	}

	return mappingIds, nil
}

func (repo *MarketplaceCmdRepo) persistInstalledItem(
	catalogItem entity.MarketplaceCatalogItem,
	hostname valueObject.Fqdn,
	installUrlPath valueObject.UrlPath,
	installDir valueObject.UnixFilePath,
	installUuid valueObject.MarketplaceInstalledItemUuid,
	mappingsId []valueObject.MappingId,
) error {
	servicesList := []string{}
	for _, service := range catalogItem.Services {
		servicesList = append(servicesList, service.String())
	}
	servicesListStr := strings.Join(servicesList, ",")

	mappingModels := []dbModel.Mapping{}
	for _, mappingId := range mappingsId {
		mappingModel := dbModel.Mapping{ID: uint(mappingId.Uint64())}
		mappingModels = append(mappingModels, mappingModel)
	}

	firstCatalogItemSlug := catalogItem.Slugs[0]
	installedItemModel := dbModel.MarketplaceInstalledItem{
		Name:             catalogItem.Name.String(),
		Hostname:         hostname.String(),
		Type:             catalogItem.Type.String(),
		UrlPath:          installUrlPath.String(),
		InstallDirectory: installDir.String(),
		InstallUuid:      installUuid.String(),
		Services:         servicesListStr,
		Mappings:         mappingModels,
		AvatarUrl:        catalogItem.AvatarUrl.String(),
		Slug:             firstCatalogItemSlug.String(),
	}

	return repo.persistentDbSvc.Handler.Create(&installedItemModel).Error
}

func (repo *MarketplaceCmdRepo) InstallItem(
	installDto dto.InstallMarketplaceCatalogItem,
) error {
	readCatalogItemDto := dto.ReadMarketplaceCatalogItemsRequest{
		Id: installDto.Id,
	}
	catalogItem, err := repo.marketplaceQueryRepo.ReadUniqueCatalogItem(
		readCatalogItemDto,
	)
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

	rawInstallUuid := uuid.New().String()[:16]
	rawInstallUuidNoHyphens := strings.Replace(rawInstallUuid, "-", "", -1)
	installUuid, err := valueObject.NewMarketplaceInstalledItemUuid(rawInstallUuidNoHyphens)
	if err != nil {
		return err
	}

	systemDataFields := repo.parseSystemDataFields(
		installDir, installUrlPath, installDto.Hostname, installUuid,
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

	err = repo.runCmdSteps("Install", catalogItem.InstallCmdSteps, receivedDataFields)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd("rm", "-rf", installTempDirPath)
	if err != nil {
		return errors.New("DeleteTmpDirectoryError: " + err.Error())
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
		catalogItem, installDto.Hostname, installUrlPath, installDir, installUuid,
		mappingIds,
	)
}

func (repo *MarketplaceCmdRepo) moveSelectedFiles(
	sourceDir valueObject.UnixFilePath,
	targetDir valueObject.UnixFilePath,
	fileNames []valueObject.UnixFileName,
	keepOnlySelectedInstead bool,
) error {
	fileNamesFilterParams := "-name \"" + fileNames[0].String() + "\""
	for _, fileToIgnore := range fileNames[1:] {
		fileNamesFilterParams += " -o -name \"" + fileToIgnore.String() + "\""
	}

	findCmdFlags := []string{"-mindepth 1", "-maxdepth 1"}
	if keepOnlySelectedInstead {
		findCmdFlags = append(findCmdFlags, "-not")
	}
	findCmdFlagsStr := strings.Join(findCmdFlags, " ")

	moveCmd := fmt.Sprintf(
		"find %s/ %s \\( %s \\) -exec mv -t %s {} +",
		sourceDir.String(), findCmdFlagsStr, fileNamesFilterParams, targetDir.String(),
	)
	_, err := infraHelper.RunCmdWithSubShell(moveCmd)
	return err
}

func (repo *MarketplaceCmdRepo) uninstallSymlinkFilesDelete(
	installedItem entity.MarketplaceInstalledItem,
	catalogItem entity.MarketplaceCatalogItem,
	softDeleteDestDirPath valueObject.UnixFilePath,
) error {
	itemHostnameStr := installedItem.Hostname.String()
	unfamiliarFilesBackupDir, err := valueObject.NewUnixFilePath(
		"/app/" + itemHostnameStr + "-unfamiliar-files-backup",
	)
	if err != nil {
		return err
	}

	unfamiliarFilesBackupDirStr := unfamiliarFilesBackupDir.String()
	err = infraHelper.MakeDir(unfamiliarFilesBackupDirStr)
	if err != nil {
		return errors.New("CreateUnfamiliarFilesBackupDirError: " + err.Error())
	}

	keepOnlySelectedInstead := true
	err = repo.moveSelectedFiles(
		installedItem.InstallDirectory, unfamiliarFilesBackupDir,
		catalogItem.UninstallFileNames, keepOnlySelectedInstead,
	)
	if err != nil {
		return errors.New("TemporarilyMoveUnfamiliarFilesError: " + err.Error())
	}

	rawInstalledItemRealRootDirPath := fmt.Sprintf(
		"/app/%s-%s-%s",
		installedItem.Slug.String(), itemHostnameStr, installedItem.InstallUuid.String(),
	)
	installedItemRealRootDirPath, err := valueObject.NewUnixFilePath(
		rawInstalledItemRealRootDirPath,
	)
	if err != nil {
		return err
	}
	installedItemRealRootDirPathStr := installedItemRealRootDirPath.String()

	_, err = infraHelper.RunCmdWithSubShell(
		"mv " + installedItemRealRootDirPathStr + "/* " + softDeleteDestDirPath.String(),
	)
	if err != nil {
		return errors.New("SoftDeleteItemFilesError: " + err.Error())
	}

	err = repo.updateFilesPrivileges(softDeleteDestDirPath)
	if err != nil {
		return errors.New("UpdateSoftDeleteDirPrivilegesError: " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"rm -rf " + installedItemRealRootDirPathStr,
	)
	if err != nil {
		return errors.New("DeleteItemRealRootPathError: " + err.Error())
	}

	itemAliasesRootDirStr := installedItem.InstallDirectory.String()
	_, err = infraHelper.RunCmdWithSubShell(
		"rm -rf " + itemAliasesRootDirStr,
	)
	if err != nil {
		return errors.New("DeleteItemAliasesRootDirError: " + err.Error())
	}

	err = infraHelper.MakeDir(itemAliasesRootDirStr)
	if err != nil {
		return errors.New("RecreateItemAliasesRootDirAsRealDirError: " + err.Error())
	}

	keepOnlySelectedInstead = true
	err = repo.moveSelectedFiles(
		unfamiliarFilesBackupDir, installedItem.InstallDirectory,
		catalogItem.UninstallFileNames, keepOnlySelectedInstead,
	)
	if err != nil {
		return errors.New("RestoreUnfamiliarFilesError: " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell("rm -rf " + unfamiliarFilesBackupDirStr)
	if err != nil {
		return errors.New("DeleteUnfamiliarFilesBackupDirError: " + err.Error())
	}

	return nil
}

func (repo *MarketplaceCmdRepo) uninstallFilesDelete(
	installedItem entity.MarketplaceInstalledItem,
	catalogItem entity.MarketplaceCatalogItem,
) error {
	if len(catalogItem.UninstallFileNames) == 0 {
		return nil
	}

	rawSoftDeleteDestDirPath := fmt.Sprintf(
		"%s/%s-%s-%s",
		useCase.TrashDirPath, installedItem.Slug.String(),
		installedItem.Hostname.String(), installedItem.InstallUuid.String(),
	)
	softDeleteDestDirPath, err := valueObject.NewUnixFilePath(rawSoftDeleteDestDirPath)
	if err != nil {
		return err
	}

	err = infraHelper.MakeDir(softDeleteDestDirPath.String())
	if err != nil {
		return errors.New("CreateSoftDeleteDirError: " + err.Error())
	}

	if infraHelper.IsSymlink(installedItem.InstallDirectory.String()) {
		return repo.uninstallSymlinkFilesDelete(
			installedItem, catalogItem, softDeleteDestDirPath,
		)
	}

	keepOnlySelectedInstead := false
	err = repo.moveSelectedFiles(
		installedItem.InstallDirectory, softDeleteDestDirPath,
		catalogItem.UninstallFileNames, keepOnlySelectedInstead,
	)
	if err != nil {
		return errors.New("SoftDeleteItemFilesError: " + err.Error())
	}

	err = repo.updateFilesPrivileges(softDeleteDestDirPath)
	if err != nil {
		return errors.New("UpdateSoftDeleteDirPrivilegesError: " + err.Error())
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

	readInstalledItemsDto := dto.ReadMarketplaceInstalledItemsRequest{
		Pagination: dto.Pagination{
			ItemsPerPage: 100,
		},
	}
	installedItemsResponseDto, err := repo.marketplaceQueryRepo.ReadInstalledItems(
		readInstalledItemsDto,
	)
	if err != nil {
		return errors.New("ReadInstalledItemsError: " + err.Error())
	}

	serviceNamesInUseMap := map[string]interface{}{}
	for _, installedItem := range installedItemsResponseDto.Items {
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

	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(repo.persistentDbSvc)
	for _, unusedService := range unusedServiceNames {
		err = servicesCmdRepo.Delete(unusedService)
		if err != nil {
			slog.Error("UninstallUnusedServiceError", slog.Any("error", err))
			continue
		}
	}

	return nil
}

func (repo *MarketplaceCmdRepo) UninstallItem(
	deleteDto dto.DeleteMarketplaceInstalledItem,
) error {
	readInstalledItemDto := dto.ReadMarketplaceInstalledItemsRequest{
		Id: &deleteDto.InstalledId,
	}
	installedItem, err := repo.marketplaceQueryRepo.ReadUniqueInstalledItem(
		readInstalledItemDto,
	)
	if err != nil {
		return err
	}

	for _, installedItemMapping := range installedItem.Mappings {
		err = repo.mappingCmdRepo.Delete(installedItemMapping.Id)
		if err != nil {
			slog.Error(
				"DeleteMappingError",
				slog.String("mappingPath", installedItemMapping.Path.String()),
				slog.Any("error", err),
			)
			continue
		}
	}

	readCatalogItemDto := dto.ReadMarketplaceCatalogItemsRequest{
		Slug: &installedItem.Slug,
	}
	catalogItem, err := repo.marketplaceQueryRepo.ReadUniqueCatalogItem(
		readCatalogItemDto,
	)
	if err != nil {
		return err
	}

	err = repo.uninstallFilesDelete(installedItem, catalogItem)
	if err != nil {
		return err
	}

	systemDataFields := repo.parseSystemDataFields(
		installedItem.InstallDirectory, installedItem.UrlPath,
		installedItem.Hostname, installedItem.InstallUuid,
	)
	err = repo.runCmdSteps("Uninstall", catalogItem.UninstallCmdSteps, systemDataFields)
	if err != nil {
		return err
	}

	installedServiceItemModel := dbModel.MarketplaceInstalledItem{
		ID: uint(deleteDto.InstalledId.Uint16()),
	}
	err = repo.persistentDbSvc.Handler.Delete(&installedServiceItemModel).Error
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

func (repo *MarketplaceCmdRepo) RefreshCatalogItems() error {
	_, err := os.Stat(infraEnvs.MarketplaceCatalogItemsDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		_, err = infraHelper.RunCmdWithSubShell(
			"cd " + infraEnvs.InfiniteOsMainDir + ";" +
				"git clone https://github.com/goinfinite/os-marketplace.git marketplace",
		)
		if err != nil {
			return errors.New("CloneMarketplaceItemsRepoError: " + err.Error())
		}
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"cd " + infraEnvs.MarketplaceCatalogItemsDir + ";" +
			"git clean -f -d; git reset --hard HEAD; git pull",
	)
	return err
}
