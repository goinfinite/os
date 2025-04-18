package marketplaceInfra

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	runtimeInfra "github.com/goinfinite/os/src/infra/runtime"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	"github.com/google/uuid"
)

const installTempDirPath = "/app/marketplace-tmp"

type MarketplaceCmdRepo struct {
	persistentDbSvc      *internalDbInfra.PersistentDatabaseService
	marketplaceQueryRepo *MarketplaceQueryRepo
	mappingCmdRepo       *vhostInfra.MappingCmdRepo
}

func NewMarketplaceCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceCmdRepo {
	return &MarketplaceCmdRepo{
		persistentDbSvc:      persistentDbSvc,
		marketplaceQueryRepo: NewMarketplaceQueryRepo(persistentDbSvc),
		mappingCmdRepo:       vhostInfra.NewMappingCmdRepo(persistentDbSvc),
	}
}

func (repo *MarketplaceCmdRepo) installServices(
	vhostName valueObject.Fqdn,
	services []valueObject.ServiceNameWithVersion,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) error {
	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(repo.persistentDbSvc)
	serviceCmdRepo := servicesInfra.NewServicesCmdRepo(repo.persistentDbSvc)

	shouldCreatePhpVirtualHost := false
	for _, serviceWithVersion := range services {
		if serviceWithVersion.Name.String() == "php-webserver" {
			shouldCreatePhpVirtualHost = true
		}

		readFirstInstalledServiceRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
			ServiceName: &serviceWithVersion.Name,
		}
		_, err := servicesQueryRepo.ReadFirstInstalledItem(
			readFirstInstalledServiceRequestDto,
		)
		if err != nil && err.Error() != servicesInfra.InstalledServiceNotFound {
			return err
		}

		if err == nil {
			continue
		}

		_, err = serviceCmdRepo.CreateInstallable(dto.CreateInstallableService{
			Name:              serviceWithVersion.Name,
			Version:           serviceWithVersion.Version,
			OperatorAccountId: operatorAccountId,
			OperatorIpAddress: operatorIpAddress,
		})
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
	itemName valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
	installDir valueObject.UnixFilePath,
	installUrlPath valueObject.UrlPath,
	installHostname valueObject.Fqdn,
	installUuid valueObject.MarketplaceInstalledItemUuid,
) (systemDataFields []valueObject.MarketplaceInstallableItemDataField) {
	dummyValueGenerator := infraHelper.DummyValueGenerator{}
	dataMap := map[string]string{
		"installDirectory":      installDir.String(),
		"installUrlPath":        installUrlPath.String(),
		"installHostname":       installHostname.String(),
		"installUuid":           installUuid.String(),
		"installTempDir":        installTempDirPath,
		"installRandomPassword": dummyValueGenerator.GenPass(16),
	}

	itemNameStr := strings.ToLower(itemName.String())
	catalogAssetsDirPath := fmt.Sprintf(
		"%s/%s/%s/assets", infraEnvs.MarketplaceCatalogItemsDir,
		itemType.String(), itemNameStr,
	)
	dataMap["marketplaceCatalogItemAssetsDirPath"] = catalogAssetsDirPath

	for key, value := range dataMap {
		dataFieldKey, _ := valueObject.NewDataFieldName(key)
		dataFieldValue, _ := valueObject.NewDataFieldValue(value)
		dataField := valueObject.NewMarketplaceInstallableItemDataField(
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

		missingDataField := valueObject.NewMarketplaceInstallableItemDataField(
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
			cmdStepStr, `%(\w{1,256})%`,
		)

		for _, cmdStepDataPlaceholder := range cmdStepDataFieldPlaceholders {
			dataFieldValue, exists := dataFieldsMap[cmdStepDataPlaceholder]
			if !exists {
				slog.Debug(
					"MissingDataField",
					slog.String("dataField", cmdStepDataPlaceholder),
				)
				dataFieldValue = ""
			}

			printableDataFieldValue := infraHelper.ShellEscape{}.StripUnsafe(dataFieldValue)

			cmdStepWithDataFieldStr := strings.ReplaceAll(
				cmdStepStr, "%"+cmdStepDataPlaceholder+"%", printableDataFieldValue,
			)
			cmdStepStr = cmdStepWithDataFieldStr
		}

		cmdStepWithDataField, _ := valueObject.NewUnixCommand(cmdStepStr)
		cmdStepsWithDataFields = append(cmdStepsWithDataFields, cmdStepWithDataField)
	}

	return cmdStepsWithDataFields, nil
}

func (repo *MarketplaceCmdRepo) runCmdSteps(
	stepsType string,
	steps []valueObject.UnixCommand,
	totalExecTimeoutSecs valueObject.UnixTime,
) error {
	if len(steps) == 0 {
		return nil
	}

	totalExecTimeoutSecsUint := uint64(totalExecTimeoutSecs.Int64())
	runCmdSettings := infraHelper.RunCmdSettings{
		ShouldRunWithSubShell: true,
		ExecutionTimeoutSecs:  totalExecTimeoutSecsUint,
	}

	totalExecRemainingTime := totalExecTimeoutSecsUint
	for stepIndex, step := range steps {
		stepStr := step.String()

		slog.Debug("Running"+stepsType+"Step", slog.String("step", stepStr))

		runCmdSettings.Command = stepStr

		stepExecTimeStart := time.Now()
		stepOutput, err := infraHelper.RunCmd(runCmdSettings)
		if err != nil {
			errorMessage := stepOutput + " | " + err.Error()
			if infraHelper.IsRunCmdTimeout(err) {
				errorMessage = "MarketplaceItem" + stepsType + "TimeoutExceeded"
			}

			return fmt.Errorf(
				"%sCmdStepError (%s): %s",
				stepsType, strconv.Itoa(stepIndex), errorMessage,
			)
		}

		stepExecElapsedTimeSecs := uint64(time.Since(stepExecTimeStart).Seconds())
		totalExecRemainingTime = totalExecRemainingTime - stepExecElapsedTimeSecs
		if totalExecRemainingTime == 0 {
			return errors.New("MarketplaceItem" + stepsType + "TimeoutExceeded")
		}

		runCmdSettings.ExecutionTimeoutSecs = totalExecRemainingTime
	}

	return nil
}

func (repo *MarketplaceCmdRepo) updateFilesPrivileges(
	targetDir valueObject.UnixFilePath,
) error {
	targetDirStr := targetDir.String()

	chownRecursively := true
	chownSymlinksToo := true
	err := infraHelper.UpdateOwnershipForWebServerUse(
		targetDirStr, chownRecursively, chownSymlinksToo,
	)
	if err != nil {
		return errors.New("ChownError (" + targetDirStr + "): " + err.Error())
	}

	dirDefaultPermissions := valueObject.NewUnixDirDefaultPermissions()
	fileDefaultPermissions := valueObject.NewUnixFileDefaultPermissions()

	updatePrivilegesCmd := fmt.Sprintf(
		"find %s -type d -exec chmod %s {} \\; && find %s -type f -exec chmod %s {} \\;",
		targetDirStr, dirDefaultPermissions.String(), targetDirStr,
		fileDefaultPermissions.String(),
	)
	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               updatePrivilegesCmd,
		ShouldRunWithSubShell: true,
	})
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
			slog.Debug(
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
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) (mappingIds []valueObject.MappingId, err error) {
	mappingQueryRepo := vhostInfra.NewMappingQueryRepo(repo.persistentDbSvc)
	mappingsReadResponse, err := mappingQueryRepo.Read(dto.ReadMappingsRequest{
		Pagination: dto.PaginationUnpaginated,
		Hostname:   &hostname,
	})
	if err != nil {
		return mappingIds, err
	}

	currentMappingsContentHashMap := map[string]entity.Mapping{}
	for _, currentMapping := range mappingsReadResponse.Mappings {
		contentHash := infraHelper.GenStrongShortHash(
			currentMapping.Hostname.String() +
				currentMapping.Path.String() +
				currentMapping.MatchPattern.String() +
				currentMapping.TargetType.String(),
		)

		currentMappingsContentHashMap[contentHash] = currentMapping
	}

	for _, itemMappingVo := range catalogMappings {
		contentHash := infraHelper.GenStrongShortHash(
			hostname.String() + itemMappingVo.Path.String() + itemMappingVo.MatchPattern.String() +
				itemMappingVo.TargetType.String(),
		)
		currentMapping, alreadyExists := currentMappingsContentHashMap[contentHash]
		if alreadyExists {
			mappingIds = append(mappingIds, currentMapping.Id)
			continue
		}

		createDto := dto.NewCreateMapping(
			hostname, itemMappingVo.Path, itemMappingVo.MatchPattern, itemMappingVo.TargetType,
			itemMappingVo.TargetValue, itemMappingVo.TargetHttpResponseCode,
			itemMappingVo.ShouldUpgradeInsecureRequests, nil, operatorAccountId, operatorIpAddress,
		)

		mappingId, err := repo.mappingCmdRepo.Create(createDto)
		if err != nil {
			slog.Debug("CreateItemMappingError", slog.String("err", err.Error()))
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
		mappingModel := dbModel.Mapping{ID: mappingId.Uint64()}
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
	if installDto.Id == nil && installDto.Slug == nil {
		return errors.New("CatalogIdOrSlugMustBeProvided")
	}

	catalogItem, err := repo.marketplaceQueryRepo.ReadFirstCatalogItem(
		dto.ReadMarketplaceCatalogItemsRequest{
			MarketplaceCatalogItemId:   installDto.Id,
			MarketplaceCatalogItemSlug: installDto.Slug,
		},
	)
	if err != nil {
		return errors.New("MarketplaceCatalogItemNotFound")
	}

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(repo.persistentDbSvc)
	vhostEntity, err := vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &installDto.Hostname,
	})
	if err != nil {
		return err
	}

	err = repo.installServices(
		installDto.Hostname, catalogItem.Services, installDto.OperatorAccountId,
		installDto.OperatorIpAddress,
	)
	if err != nil {
		return err
	}

	installUrlPath, _ := valueObject.NewUrlPath("/")
	if installDto.UrlPath != nil {
		installUrlPath = *installDto.UrlPath
	}

	installDirStr := vhostEntity.RootDirectory.String() + installUrlPath.GetWithoutTrailingSlash()
	installDir, err := valueObject.NewUnixFilePath(installDirStr)
	if err != nil {
		return errors.New("DefineInstallDirectoryError: " + err.Error())
	}
	installDirStr = installDir.String()
	if installDirStr == infraEnvs.PrimaryPublicDir {
		err := infraHelper.BackupPrimaryIndexFile()
		if err != nil {
			return err
		}
	}

	rawInstallUuid := uuid.New().String()[:16]
	rawInstallUuidNoHyphens := strings.Replace(rawInstallUuid, "-", "", -1)
	installUuid, err := valueObject.NewMarketplaceInstalledItemUuid(rawInstallUuidNoHyphens)
	if err != nil {
		return err
	}

	systemDataFields := repo.parseSystemDataFields(
		catalogItem.Name, catalogItem.Type, installDir, installUrlPath,
		installDto.Hostname, installUuid,
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

	usableInstallCmdSteps, err := repo.replaceCmdStepsPlaceholders(
		catalogItem.InstallCmdSteps, receivedDataFields,
	)
	if err != nil {
		return errors.New("ParseCmdStepWithDataFieldsError: " + err.Error())
	}

	err = repo.runCmdSteps(
		"Install", usableInstallCmdSteps, catalogItem.InstallTimeoutSecs,
	)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "rm",
		Args:    []string{"-rf", installTempDirPath},
	})
	if err != nil {
		return errors.New("DeleteTmpDirectoryError: " + err.Error())
	}

	err = repo.updateFilesPrivileges(installDir)
	if err != nil {
		return errors.New("UpdateFilesPrivilegesError: " + err.Error())
	}

	isRootDirectory := installDir.String() == vhostEntity.RootDirectory.String()
	if !isRootDirectory {
		catalogItem.Mappings = repo.updateMappingsBase(
			catalogItem.Mappings, installUrlPath,
		)
	}

	mappingIds, err := repo.createMappings(
		installDto.Hostname, catalogItem.Mappings, installDto.OperatorAccountId,
		installDto.OperatorIpAddress,
	)
	if err != nil {
		return err
	}

	err = repo.persistInstalledItem(
		catalogItem, installDto.Hostname, installUrlPath, installDir, installUuid,
		mappingIds,
	)
	if err != nil {
		return errors.New("PersistInstalledItemError: " + err.Error())
	}

	for _, mappingId := range mappingIds {
		err = repo.mappingCmdRepo.UpdateMarketplaceItem(mappingId, catalogItem.Name)
		if err != nil {
			slog.Debug(
				"UpdateMappingItemNameError",
				slog.String("mappingId", mappingId.String()),
				slog.String("err", err.Error()),
			)
			continue
		}
	}

	return nil
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
	_, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               moveCmd,
		ShouldRunWithSubShell: true,
	})
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

	softDeleteCmd := "mv " + installedItemRealRootDirPathStr + "/* " + softDeleteDestDirPath.String()
	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               softDeleteCmd,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return errors.New("SoftDeleteItemFilesError: " + err.Error())
	}

	err = repo.updateFilesPrivileges(softDeleteDestDirPath)
	if err != nil {
		return errors.New("UpdateSoftDeleteDirPrivilegesError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               "rm -rf " + installedItemRealRootDirPathStr,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return errors.New("DeleteItemRealRootPathError: " + err.Error())
	}

	itemAliasesRootDirStr := installedItem.InstallDirectory.String()
	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               "rm -rf " + itemAliasesRootDirStr,
		ShouldRunWithSubShell: true,
	})
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

	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               "rm -rf " + unfamiliarFilesBackupDirStr,
		ShouldRunWithSubShell: true,
	})
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
		valueObject.UnixFilePathTrashDir.String(), installedItem.Slug.String(),
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

	if installedItem.InstallDirectory.String() == infraEnvs.PrimaryPublicDir {
		err := infraHelper.RestorePrimaryIndexFile()
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *MarketplaceCmdRepo) UninstallItem(
	deleteDto dto.DeleteMarketplaceInstalledItem,
) error {
	installedItem, err := repo.marketplaceQueryRepo.ReadFirstInstalledItem(
		dto.ReadMarketplaceInstalledItemsRequest{
			MarketplaceInstalledItemId: &deleteDto.InstalledId,
		},
	)
	if err != nil {
		return errors.New("ReadMarketplaceInstalledItemError: " + err.Error())
	}

	catalogItem, err := repo.marketplaceQueryRepo.ReadFirstCatalogItem(
		dto.ReadMarketplaceCatalogItemsRequest{
			MarketplaceCatalogItemSlug: &installedItem.Slug,
		},
	)
	if err != nil {
		return errors.New("ReadMarketplaceCatalogItemError: " + err.Error())
	}

	err = repo.uninstallFilesDelete(installedItem, catalogItem)
	if err != nil {
		return err
	}

	systemDataFields := repo.parseSystemDataFields(
		installedItem.Name, installedItem.Type, installedItem.InstallDirectory,
		installedItem.UrlPath, installedItem.Hostname, installedItem.InstallUuid,
	)
	usableInstallCmdSteps, err := repo.replaceCmdStepsPlaceholders(
		catalogItem.UninstallCmdSteps, systemDataFields,
	)
	if err != nil {
		return errors.New("ParseCmdStepWithDataFieldsError: " + err.Error())
	}

	err = repo.runCmdSteps(
		"Uninstall", usableInstallCmdSteps, catalogItem.UninstallTimeoutSecs,
	)
	if err != nil {
		return err
	}

	marketplaceInstalledItemModel := dbModel.MarketplaceInstalledItem{
		ID: deleteDto.InstalledId.Uint16(),
	}
	return repo.persistentDbSvc.Handler.Delete(&marketplaceInstalledItemModel).Error
}

func (repo *MarketplaceCmdRepo) RefreshCatalogItems() error {
	_, err := os.Stat(infraEnvs.MarketplaceCatalogItemsDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		repoCloneCmd := fmt.Sprintf(
			"cd %s; git clone --single-branch --branch %s %s marketplace",
			infraEnvs.InfiniteOsMainDir, infraEnvs.MarketplaceCatalogItemsRepoBranch,
			infraEnvs.MarketplaceCatalogItemsRepoUrl,
		)
		_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
			Command:               repoCloneCmd,
			ShouldRunWithSubShell: true,
		})
		if err != nil {
			return errors.New("CloneMarketplaceItemsRepoError: " + err.Error())
		}
	}

	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "cd " + infraEnvs.MarketplaceCatalogItemsDir + ";" +
			"git clean -f -d; git reset --hard HEAD; git pull",
		ShouldRunWithSubShell: true,
	})
	return err
}
