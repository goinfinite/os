package servicesInfra

import (
	"errors"
	"log/slog"
	"math"
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
	"github.com/iancoleman/strcase"

	"github.com/shirou/gopsutil/process"
)

type ServicesQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewServicesQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ServicesQueryRepo {
	return &ServicesQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ServicesQueryRepo) readPidProcessFamily(pid int32) ([]*process.Process, error) {
	processFamily := []*process.Process{}

	pidProcess, err := process.NewProcess(pid)
	if err != nil {
		return processFamily, err
	}

	processFamily = append(processFamily, pidProcess)

	childrenPidProcesses, err := pidProcess.Children()
	if err != nil || len(childrenPidProcesses) == 0 {
		return processFamily, nil
	}

	for _, childPidProcess := range childrenPidProcesses {
		grandChildrenPidProcesses, err := repo.readPidProcessFamily(
			childPidProcess.Pid,
		)
		if err != nil || len(grandChildrenPidProcesses) == 0 {
			continue
		}

		processFamily = append(processFamily, grandChildrenPidProcesses...)
	}

	return processFamily, nil
}

func (repo *ServicesQueryRepo) readPidMetrics(
	mainPid int32,
) (serviceMetrics valueObject.ServiceMetrics, err error) {
	pidProcesses, err := repo.readPidProcessFamily(mainPid)
	if err != nil {
		return serviceMetrics, err
	}

	if len(pidProcesses) == 0 {
		return serviceMetrics, nil
	}

	uptimeMilliseconds, err := pidProcesses[0].CreateTime()
	if err != nil {
		return serviceMetrics, err
	}
	nowMilliseconds := time.Now().UTC().UnixMilli()
	uptimeSecs := (nowMilliseconds - uptimeMilliseconds) / 1000

	cpuPercent := float64(0.0)
	memPercent := float32(0.0)

	pids := []uint32{}
	for _, process := range pidProcesses {
		pidCpuPercent, err := process.CPUPercent()
		if err != nil {
			slog.Debug(err.Error(), slog.Int("processPid", int(process.Pid)))
			continue
		}

		pidMemPercent, err := process.MemoryPercent()
		if err != nil {
			slog.Debug(err.Error(), slog.Int("processPid", int(process.Pid)))
			continue
		}

		cpuPercent += pidCpuPercent
		memPercent += pidMemPercent

		pids = append(pids, uint32(process.Pid))
	}

	cpuPercent = math.Round(cpuPercent*100) / 100
	memPercent = float32(math.Round(float64(memPercent)*100) / 100)

	serviceMetrics = valueObject.NewServiceMetrics(
		pids, uptimeSecs, cpuPercent, memPercent,
	)

	return serviceMetrics, nil
}

func (repo *ServicesQueryRepo) readServiceMetrics(
	name valueObject.ServiceName,
) (*valueObject.ServiceMetrics, error) {
	supervisorStatus, _ := infraHelper.RunCmdWithSubShell(
		SupervisorCtlBin + " status " + name.String(),
	)
	if len(supervisorStatus) == 0 {
		return nil, errors.New("ReadSupervisorStatusError")
	}

	// # supervisorctl status <serviceName>
	// <serviceName>                    RUNNING   pid 120, uptime 0:00:35
	supervisorStatusParts := strings.Fields(supervisorStatus)
	if len(supervisorStatusParts) < 4 {
		return nil, errors.New("MissingSupervisorStatusParts")
	}

	rawServiceStatus := supervisorStatusParts[1]
	serviceStatus, err := valueObject.NewServiceStatus(rawServiceStatus)
	if err != nil {
		return nil, errors.New(err.Error() + ": " + rawServiceStatus)
	}

	if serviceStatus.String() != "running" {
		return nil, nil
	}

	rawServicePid := supervisorStatusParts[3]
	rawServicePid = strings.Trim(rawServicePid, ",")
	servicePidInt, err := strconv.ParseInt(rawServicePid, 10, 32)
	if err != nil {
		return nil, errors.New(err.Error() + ": " + rawServicePid)
	}

	serviceMetrics, err := repo.readPidMetrics(int32(servicePidInt))
	if err != nil {
		return nil, errors.New(err.Error() + ": " + rawServicePid)
	}

	return &serviceMetrics, nil
}

func (repo *ServicesQueryRepo) readStoppedServicesNames() ([]string, error) {
	stoppedServicesNames := []string{}

	rawStoppedServices, err := infraHelper.RunCmdWithSubShell(
		SupervisorCtlBin + " status | grep -v 'RUNNING' | awk '{print $1}'",
	)
	if err != nil {
		return stoppedServicesNames, err
	}

	rawStoppedServicesLines := strings.Split(rawStoppedServices, "\n")
	for _, rawStoppedService := range rawStoppedServicesLines {
		if rawStoppedService == "" {
			continue
		}

		serviceName, err := valueObject.NewServiceName(rawStoppedService)
		if err != nil {
			slog.Debug(
				"InvalidStoppedServiceName",
				slog.String("serviceName", rawStoppedService),
			)
			continue
		}

		stoppedServicesNames = append(stoppedServicesNames, serviceName.String())
	}

	return stoppedServicesNames, nil
}

func (repo *ServicesQueryRepo) ReadInstalledItems(
	readDto dto.ReadInstalledServicesItemsRequest,
) (installedItemsDto dto.ReadInstalledServicesItemsResponse, err error) {
	model := dbModel.InstalledService{}
	if readDto.Name != nil {
		model.Name = readDto.Name.String()
	}
	if readDto.Nature != nil {
		model.Nature = readDto.Nature.String()
	}
	if readDto.Type != nil {
		model.Type = readDto.Type.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.
		Where(&model).
		Limit(int(readDto.Pagination.ItemsPerPage))
	if readDto.Pagination.LastSeenId == nil {
		offset := int(readDto.Pagination.PageNumber) * int(readDto.Pagination.ItemsPerPage)
		dbQuery = dbQuery.Offset(offset)
	} else {
		dbQuery = dbQuery.Where("id > ?", readDto.Pagination.LastSeenId.String())
	}
	if readDto.Pagination.SortBy != nil {
		orderStatement := readDto.Pagination.SortBy.String()
		orderStatement = strcase.ToSnake(orderStatement)
		if orderStatement == "id" {
			orderStatement = "ID"
		}

		if readDto.Pagination.SortDirection != nil {
			orderStatement += " " + readDto.Pagination.SortDirection.String()
		}

		dbQuery = dbQuery.Order(orderStatement)
	}

	models := []dbModel.InstalledService{}
	err = dbQuery.Find(&models).Error
	if err != nil {
		return installedItemsDto, errors.New("ReadInstalledServicesItemsError")
	}

	var itemsTotal int64
	err = dbQuery.Count(&itemsTotal).Error
	if err != nil {
		return installedItemsDto, errors.New(
			"CountInstalledServicesItemsTotalError: " + err.Error(),
		)
	}

	entities := []dto.InstalledServiceWithMetrics{}
	for _, model := range models {
		entityWithoutMetrics, err := model.ToEntity()
		if err != nil {
			slog.Error(
				"InstalledServiceItemModelToEntityError",
				slog.String("name", model.Name), slog.Any("error", err),
			)
			continue
		}

		var entityMetricsPtr *valueObject.ServiceMetrics
		if readDto.ShouldIncludeMetrics {
			entityMetricsPtr, err = repo.readServiceMetrics(entityWithoutMetrics.Name)
			if err != nil {
				slog.Error(
					"FailedToReadInstalledServiceMetrics",
					slog.String("name", model.Name), slog.Any("error", err),
				)
				entityMetricsPtr = nil
			}
		}

		entityWithMetrics := dto.NewInstalledServiceWithMetrics(
			entityWithoutMetrics,
			entityMetricsPtr,
		)
		entities = append(entities, entityWithMetrics)
	}

	stoppedServicesNames, err := repo.readStoppedServicesNames()
	if err != nil {
		return installedItemsDto, errors.New(
			"FailedToReadStoppedServicesNames: " + err.Error(),
		)
	}

	stoppedStatus, _ := valueObject.NewServiceStatus("stopped")
	for entityIndex, entity := range entities {
		if !slices.Contains(stoppedServicesNames, entity.Name.String()) {
			continue
		}

		entities[entityIndex].Status = stoppedStatus
	}

	itemsTotalUint := uint64(itemsTotal)
	pagesTotal := uint32(
		math.Ceil(float64(itemsTotal) / float64(readDto.Pagination.ItemsPerPage)),
	)
	responsePagination := dto.Pagination{
		PageNumber:    readDto.Pagination.PageNumber,
		ItemsPerPage:  readDto.Pagination.ItemsPerPage,
		SortBy:        readDto.Pagination.SortBy,
		SortDirection: readDto.Pagination.SortDirection,
		PagesTotal:    &pagesTotal,
		ItemsTotal:    &itemsTotalUint,
	}

	return dto.ReadInstalledServicesItemsResponse{
		Pagination: responsePagination,
		Items:      entities,
	}, nil
}

func (repo *ServicesQueryRepo) ReadUniqueInstalledItem(
	readDto dto.ReadInstalledServicesItemsRequest,
) (installedItem dto.InstalledServiceWithMetrics, err error) {
	readDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.ReadInstalledItems(readDto)
	if err != nil {
		return installedItem, err
	}

	if len(responseDto.Items) == 0 {
		return installedItem, errors.New("ServiceInstalledItemNotFound")
	}

	foundInstalledItem := responseDto.Items[0]
	return foundInstalledItem, nil
}

func (repo *ServicesQueryRepo) parseManifestCmdSteps(
	stepsType string,
	rawCmdSteps interface{},
) (cmdSteps []valueObject.UnixCommand, err error) {
	cmdStepsMap, assertOk := rawCmdSteps.([]interface{})
	if !assertOk {
		return cmdSteps, errors.New("InvalidCmdStepsStructure")
	}

	for _, rawCmd := range cmdStepsMap {
		command, err := valueObject.NewUnixCommand(rawCmd)
		if err != nil {
			slog.Debug(
				"ParseInvalidCmdStepError",
				slog.String("stepsType", stepsType),
				slog.Any("rawCmd", rawCmd),
			)
			return cmdSteps, err
		}
		cmdSteps = append(cmdSteps, command)
	}

	return cmdSteps, nil
}

func (repo *ServicesQueryRepo) installableServiceFactory(
	serviceFilePath valueObject.UnixFilePath,
) (installableService entity.InstallableService, err error) {
	serviceMap, err := infraHelper.FileSerializedDataToMap(serviceFilePath)
	if err != nil {
		return installableService, err
	}

	requiredParams := []string{
		"name", "nature", "type", "startCmd", "description", "installCmdSteps",
	}
	for _, requiredParam := range requiredParams {
		if serviceMap[requiredParam] != nil {
			continue
		}

		return installableService, errors.New("MissingParam: " + requiredParam)
	}

	name, err := valueObject.NewServiceName(serviceMap["name"])
	if err != nil {
		return installableService, err
	}
	nameStr := name.String()

	nature, err := valueObject.NewServiceNature(serviceMap["nature"])
	if err != nil {
		return installableService, err
	}

	serviceType, err := valueObject.NewServiceType(serviceMap["type"])
	if err != nil {
		return installableService, err
	}

	startCommand, err := valueObject.NewUnixCommand(serviceMap["startCmd"])
	if err != nil {
		return installableService, err
	}

	description, err := valueObject.NewServiceDescription(serviceMap["description"])
	if err != nil {
		return installableService, err
	}

	versions := []valueObject.ServiceVersion{}
	if serviceMap["versions"] != nil {
		versionsMap, assertOk := serviceMap["versions"].([]interface{})
		if !assertOk {
			return installableService, errors.New("InvalidServiceVersions")
		}
		for _, rawVersion := range versionsMap {
			version, err := valueObject.NewServiceVersion(rawVersion)
			if err != nil {
				slog.Debug(
					"ParseInvalidServiceVersionError",
					slog.String("serviceName", nameStr),
					slog.Any("version", rawVersion),
				)
				continue
			}
			versions = append(versions, version)
		}
	}

	envs := []valueObject.ServiceEnv{}
	if serviceMap["envs"] != nil {
		envsMap, assertOk := serviceMap["envs"].([]interface{})
		if !assertOk {
			return installableService, errors.New("InvalidEnvs")
		}
		for _, rawEnv := range envsMap {
			env, err := valueObject.NewServiceEnv(rawEnv)
			if err != nil {
				slog.Debug(
					"ParseInvalidEnvError",
					slog.String("serviceName", nameStr),
					slog.Any("env", rawEnv),
				)
				continue
			}
			envs = append(envs, env)
		}
	}

	portBindings := []valueObject.PortBinding{}
	if serviceMap["portBindings"] != nil {
		portBindingsMap, assertOk := serviceMap["portBindings"].([]interface{})
		if !assertOk {
			return installableService, errors.New("InvalidPortBindings")
		}
		for _, rawPortBinding := range portBindingsMap {
			portBinding, err := valueObject.NewPortBinding(rawPortBinding)
			if err != nil {
				slog.Debug(
					"ParseInvalidPortBindingError",
					slog.String("serviceName", nameStr),
					slog.Any("portBinding", rawPortBinding),
				)
				continue
			}
			portBindings = append(portBindings, portBinding)
		}
	}

	stopCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["stopCmdSteps"] != nil {
		stopCmdSteps, err = repo.parseManifestCmdSteps(
			"Stop", serviceMap["stopCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	installCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["installCmdSteps"] != nil {
		installCmdSteps, err = repo.parseManifestCmdSteps(
			"Install", serviceMap["installCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	uninstallCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["uninstallCmdSteps"] != nil {
		uninstallCmdSteps, err = repo.parseManifestCmdSteps(
			"Uninstall", serviceMap["uninstallCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	uninstallFilePaths := []valueObject.UnixFilePath{}
	if serviceMap["uninstallFilePaths"] != nil {
		filesMap, assertOk := serviceMap["uninstallFilePaths"].([]interface{})
		if !assertOk {
			return installableService, errors.New("InvalidUninstallFilePaths")
		}
		for _, rawFileName := range filesMap {
			fileName, err := valueObject.NewUnixFilePath(rawFileName)
			if err != nil {
				slog.Debug(
					"ParseInvalidUninstallFilePathError",
					slog.String("serviceName", nameStr),
					slog.Any("filePath", rawFileName),
				)
				continue
			}
			uninstallFilePaths = append(uninstallFilePaths, fileName)
		}
	}

	preStartCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["preStartCmdSteps"] != nil {
		preStartCmdSteps, err = repo.parseManifestCmdSteps(
			"PreStart", serviceMap["preStartCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	postStartCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["postStartCmdSteps"] != nil {
		postStartCmdSteps, err = repo.parseManifestCmdSteps(
			"PostStart", serviceMap["postStartCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	preStopCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["preStopCmdSteps"] != nil {
		preStopCmdSteps, err = repo.parseManifestCmdSteps(
			"PreStop", serviceMap["preStopCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	postStopCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["postStopCmdSteps"] != nil {
		postStopCmdSteps, err = repo.parseManifestCmdSteps(
			"PostStop", serviceMap["postStopCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	var execUserPtr *valueObject.UnixUsername
	if serviceMap["execUser"] != nil {
		execUser, err := valueObject.NewUnixUsername(serviceMap["execUser"])
		if err != nil {
			return installableService, err
		}
		execUserPtr = &execUser
	}

	var workingDirectoryPtr *valueObject.UnixFilePath
	if serviceMap["workingDirectory"] != nil {
		workingDirectory, err := valueObject.NewUnixFilePath(serviceMap["workingDirectory"])
		if err != nil {
			return installableService, err
		}
		workingDirectoryPtr = &workingDirectory
	}

	var startupFilePtr *valueObject.UnixFilePath
	if serviceMap["startupFile"] != nil {
		startupFile, err := valueObject.NewUnixFilePath(serviceMap["startupFile"])
		if err != nil {
			return installableService, err
		}
		startupFilePtr = &startupFile
	}

	var logOutputPathPtr *valueObject.UnixFilePath
	if serviceMap["logOutputPath"] != nil {
		logOutputPath, err := valueObject.NewUnixFilePath(serviceMap["logOutputPath"])
		if err != nil {
			return installableService, err
		}
		logOutputPathPtr = &logOutputPath
	}

	var logErrorPathPtr *valueObject.UnixFilePath
	if serviceMap["logErrorPath"] != nil {
		logErrorPath, err := valueObject.NewUnixFilePath(serviceMap["logErrorPath"])
		if err != nil {
			return installableService, err
		}
		logErrorPathPtr = &logErrorPath
	}

	var estimatedSizeBytesPtr *valueObject.Byte
	if serviceMap["estimatedSizeBytes"] != nil {
		estimatedSizeBytes, err := valueObject.NewByte(serviceMap["estimatedSizeBytes"])
		if err != nil {
			return installableService, err
		}
		estimatedSizeBytesPtr = &estimatedSizeBytes
	}

	var avatarUrlPtr *valueObject.Url
	if serviceMap["avatarUrl"] != nil {
		avatarUrl, err := valueObject.NewUrl(serviceMap["avatarUrl"])
		if err != nil {
			return installableService, err
		}
		avatarUrlPtr = &avatarUrl
	}

	return entity.NewInstallableService(
		name, nature, serviceType, startCommand, description, versions, envs,
		portBindings, stopCmdSteps, installCmdSteps, uninstallCmdSteps, uninstallFilePaths,
		preStartCmdSteps, postStartCmdSteps, preStopCmdSteps, postStopCmdSteps,
		execUserPtr, workingDirectoryPtr, startupFilePtr, logOutputPathPtr,
		logErrorPathPtr, estimatedSizeBytesPtr, avatarUrlPtr,
	), nil
}

func (repo *ServicesQueryRepo) ReadInstallableItems(
	readDto dto.ReadInstallableServicesItemsRequest,
) (installableItemsDto dto.ReadInstallableServicesItemsResponse, err error) {
	_, err = os.Stat(infraEnvs.ServicesItemsDir)
	if err != nil {
		servicesCmdRepo := NewServicesCmdRepo(repo.persistentDbSvc)
		err = servicesCmdRepo.RefreshItems()
		if err != nil {
			return installableItemsDto, errors.New(
				"RefreshServicesItemsError: " + err.Error(),
			)
		}
	}

	rawInstallableFilesList, err := infraHelper.RunCmdWithSubShell(
		"find " + infraEnvs.ServicesItemsDir + " -type f " +
			"\\( -name '*.json' -o -name '*.yaml' -o -name '*.yml' \\) " +
			"-not -path '*/.*' -not -name '.*'",
	)
	if err != nil {
		return installableItemsDto, errors.New(
			"ReadInstallableFilesError: " + err.Error(),
		)
	}

	if len(rawInstallableFilesList) == 0 {
		return installableItemsDto, errors.New("NoInstallableFilesFound")
	}

	rawInstallableFilesListParts := strings.Split(rawInstallableFilesList, "\n")
	if len(rawInstallableFilesListParts) == 0 {
		return installableItemsDto, errors.New("NoInstallableFilesFound")
	}

	installableServices := []entity.InstallableService{}
	for _, rawFilePath := range rawInstallableFilesListParts {
		itemFilePath, err := valueObject.NewUnixFilePath(rawFilePath)
		if err != nil {
			slog.Error(err.Error(), slog.String("filePath", rawFilePath))
			continue
		}

		installableService, err := repo.installableServiceFactory(itemFilePath)
		if err != nil {
			slog.Error(
				"CatalogMarketplaceItemFactoryError",
				slog.String("filePath", itemFilePath.String()), slog.Any("err", err),
			)
			continue
		}

		if readDto.Name != nil {
			isNameEqual := strings.EqualFold(
				installableService.Name.String(), readDto.Name.String(),
			)
			if !isNameEqual {
				continue
			}
		}

		if readDto.Nature != nil {
			isNatureEqual := strings.EqualFold(
				installableService.Nature.String(), readDto.Nature.String(),
			)
			if !isNatureEqual {
				continue
			}
		}

		if readDto.Type != nil && installableService.Type != *readDto.Type {
			isTypeEqual := strings.EqualFold(
				installableService.Type.String(), readDto.Type.String(),
			)
			if !isTypeEqual {
				continue
			}
		}

		installableServices = append(installableServices, installableService)
	}

	sortDirectionStr := "asc"
	if readDto.Pagination.SortDirection != nil {
		sortDirectionStr = readDto.Pagination.SortDirection.String()
	}

	if readDto.Pagination.SortBy != nil {
		slices.SortStableFunc(installableServices, func(a, b entity.InstallableService) int {
			firstElement := a
			secondElement := b
			if sortDirectionStr != "asc" {
				firstElement = b
				secondElement = a
			}

			switch readDto.Pagination.SortBy.String() {
			case "name":
				return strings.Compare(
					firstElement.Name.String(), secondElement.Name.String(),
				)
			case "nature":
				return strings.Compare(
					firstElement.Nature.String(), secondElement.Nature.String(),
				)
			case "type":
				return strings.Compare(
					firstElement.Type.String(), secondElement.Type.String(),
				)
			default:
				return 0
			}
		})
	}

	itemsTotal := uint64(len(installableServices))
	pagesTotal := uint32(itemsTotal / uint64(readDto.Pagination.ItemsPerPage))

	paginationDto := readDto.Pagination
	paginationDto.ItemsTotal = &itemsTotal
	paginationDto.PagesTotal = &pagesTotal

	return dto.ReadInstallableServicesItemsResponse{
		Pagination: paginationDto,
		Items:      installableServices,
	}, nil
}

func (repo *ServicesQueryRepo) ReadUniqueInstallableItem(
	readDto dto.ReadInstallableServicesItemsRequest,
) (installableService entity.InstallableService, err error) {
	readDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.ReadInstallableItems(readDto)
	if err != nil {
		return installableService, err
	}

	if len(responseDto.Items) == 0 {
		return installableService, errors.New("InstallableServiceItemNotFound")
	}

	foundInstallableItem := responseDto.Items[0]
	return foundInstallableItem, nil
}
