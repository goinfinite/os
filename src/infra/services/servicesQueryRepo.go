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
	dbHelper "github.com/goinfinite/os/src/infra/internalDatabase/helper"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	tkInfra "github.com/goinfinite/tk/src/infra"

	"github.com/shirou/gopsutil/process"
)

const InstalledServiceNotFound = "ServiceInstalledItemNotFound"

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

func (repo *ServicesQueryRepo) readStoppedServicesNames() ([]string, error) {
	stoppedServicesNames := []string{}

	readStoppedServicesCmd := SupervisorCtlBin + " status | grep -v 'RUNNING' | awk '{print $1}'"
	rawStoppedServices, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               readStoppedServicesCmd,
		ShouldRunWithSubShell: true,
	})
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

func (repo *ServicesQueryRepo) installedServicesMetricsFactory(
	installedServices []entity.InstalledService,
) []dto.InstalledServiceWithMetrics {
	installedServicesWithMetrics := []dto.InstalledServiceWithMetrics{}
	for _, installedService := range installedServices {
		serviceWithoutMetrics := dto.NewInstalledServiceWithMetrics(
			installedService, nil,
		)

		serviceNameStr := installedService.Name.String()

		supervisorStatus, _ := infraHelper.RunCmd(infraHelper.RunCmdSettings{
			Command:               SupervisorCtlBin + " status " + serviceNameStr,
			ShouldRunWithSubShell: true,
		})
		if len(supervisorStatus) == 0 {
			installedServicesWithMetrics = append(
				installedServicesWithMetrics, serviceWithoutMetrics,
			)

			slog.Debug("ReadSupervisorStatusError", slog.String("name", serviceNameStr))
			continue
		}

		// # supervisorctl status <serviceName>
		// <serviceName>                    RUNNING   pid 120, uptime 0:00:35
		supervisorStatusParts := strings.Fields(supervisorStatus)
		if len(supervisorStatusParts) < 4 {
			slog.Debug("MissingSupervisorStatusParts", slog.String("name", serviceNameStr))
		}

		rawServiceStatus := supervisorStatusParts[1]
		serviceStatus, err := valueObject.NewServiceStatus(rawServiceStatus)
		if err != nil {
			installedServicesWithMetrics = append(
				installedServicesWithMetrics, serviceWithoutMetrics,
			)

			slog.Debug(
				err.Error(), slog.String("name", serviceNameStr),
				slog.String("rawStatus", rawServiceStatus),
			)
			continue
		}

		if serviceStatus.String() != "running" {
			installedServicesWithMetrics = append(
				installedServicesWithMetrics, serviceWithoutMetrics,
			)

			continue
		}

		rawServicePid := supervisorStatusParts[3]
		rawServicePid = strings.Trim(rawServicePid, ",")
		servicePidInt, err := strconv.ParseInt(rawServicePid, 10, 32)
		if err != nil {
			installedServicesWithMetrics = append(
				installedServicesWithMetrics, serviceWithoutMetrics,
			)

			slog.Debug(
				err.Error(), slog.String("name", serviceNameStr),
				slog.String("rawPid", rawServicePid),
			)
			continue
		}

		metrics, err := repo.readPidMetrics(int32(servicePidInt))
		if err != nil {
			installedServicesWithMetrics = append(
				installedServicesWithMetrics, serviceWithoutMetrics,
			)

			slog.Debug(
				err.Error(), slog.String("name", serviceNameStr),
				slog.String("rawPid", rawServicePid),
			)
			continue
		}

		serviceWithMetrics := dto.NewInstalledServiceWithMetrics(
			installedService, &metrics,
		)
		installedServicesWithMetrics = append(
			installedServicesWithMetrics, serviceWithMetrics,
		)
	}

	return installedServicesWithMetrics
}

func (repo *ServicesQueryRepo) ReadInstalledItems(
	requestDto dto.ReadInstalledServicesItemsRequest,
) (installedItemsDto dto.ReadInstalledServicesItemsResponse, err error) {
	installedServiceModel := dbModel.InstalledService{}
	if requestDto.ServiceNature != nil {
		installedServiceModel.Nature = requestDto.ServiceNature.String()
	}
	if requestDto.ServiceType != nil {
		installedServiceModel.Type = requestDto.ServiceType.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.
		Model(&installedServiceModel).
		Where(&installedServiceModel)
	if requestDto.ServiceName != nil {
		serviceNameLike := "%" + requestDto.ServiceName.String() + "%"
		dbQuery = dbQuery.Where("name LIKE ?", serviceNameLike)
	}

	paginatedDbQuery, responsePagination, err := dbHelper.PaginationQueryBuilder(
		dbQuery, requestDto.Pagination,
	)
	if err != nil {
		return installedItemsDto, errors.New(
			"PaginationQueryBuilderError: " + err.Error(),
		)
	}

	installedServiceModels := []dbModel.InstalledService{}
	err = paginatedDbQuery.Find(&installedServiceModels).Error
	if err != nil {
		return installedItemsDto, errors.New("ReadInstalledServicesItemsError")
	}

	installedServiceEntities := []entity.InstalledService{}
	for _, resultModel := range installedServiceModels {
		entity, err := resultModel.ToEntity()
		if err != nil {
			slog.Debug(
				"InstalledServiceItemModelToEntityError",
				slog.String("name", resultModel.Name),
				slog.String("err", err.Error()),
			)
			continue
		}

		installedServiceEntities = append(installedServiceEntities, entity)
	}

	stoppedServicesNames, err := repo.readStoppedServicesNames()
	if err != nil {
		return installedItemsDto, errors.New(
			"FailedToReadStoppedServicesNames: " + err.Error(),
		)
	}

	stoppedStatus, _ := valueObject.NewServiceStatus("stopped")
	for serviceEntityIndex, serviceEntity := range installedServiceEntities {
		if !slices.Contains(stoppedServicesNames, serviceEntity.Name.String()) {
			continue
		}

		installedServiceEntities[serviceEntityIndex].Status = stoppedStatus
	}

	if requestDto.ServiceStatus != nil {
		filteredServiceEntities := []entity.InstalledService{}
		for _, serviceEntity := range installedServiceEntities {
			if serviceEntity.Status != *requestDto.ServiceStatus {
				continue
			}

			filteredServiceEntities = append(filteredServiceEntities, serviceEntity)
		}

		installedServiceEntities = filteredServiceEntities

		itemsTotal := uint64(len(filteredServiceEntities))
		responsePagination.ItemsTotal = &itemsTotal

		pagesTotal := uint32(
			math.Ceil(float64(itemsTotal) / float64(responsePagination.ItemsPerPage)),
		)
		responsePagination.PagesTotal = &pagesTotal
	}
	responseDto := dto.ReadInstalledServicesItemsResponse{
		Pagination: responsePagination,
	}

	if requestDto.ShouldIncludeMetrics != nil && *requestDto.ShouldIncludeMetrics {
		responseDto.InstalledServicesWithMetrics = repo.installedServicesMetricsFactory(
			installedServiceEntities,
		)
		return responseDto, nil
	}

	responseDto.InstalledServices = installedServiceEntities
	return responseDto, nil
}

func (repo *ServicesQueryRepo) ReadFirstInstalledItem(
	readFirstRequestDto dto.ReadFirstInstalledServiceItemsRequest,
) (installedItem entity.InstalledService, err error) {
	shouldIncludeMetrics := false
	readRequestDto := dto.ReadInstalledServicesItemsRequest{
		Pagination: dto.Pagination{
			PageNumber:   0,
			ItemsPerPage: 1,
		},
		ServiceName:          readFirstRequestDto.ServiceName,
		ServiceNature:        readFirstRequestDto.ServiceNature,
		ServiceType:          readFirstRequestDto.ServiceType,
		ShouldIncludeMetrics: &shouldIncludeMetrics,
	}
	responseDto, err := repo.ReadInstalledItems(readRequestDto)
	if err != nil {
		return installedItem, err
	}

	if len(responseDto.InstalledServices) == 0 {
		return installedItem, errors.New(InstalledServiceNotFound)
	}

	return responseDto.InstalledServices[0], nil
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
	serviceMap, err := tkInfra.FileDeserializer(serviceFilePath.String())
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

	manifestVersion, _ := valueObject.NewServiceManifestVersion("v1")
	if serviceMap["manifestVersion"] != nil {
		manifestVersion, err = valueObject.NewServiceManifestVersion(
			serviceMap["manifestVersion"],
		)
		if err != nil {
			return installableService, err
		}
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
			return installableService, errors.New("InvalidServiceVersionsStructure")
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
			return installableService, errors.New("InvalidPortBindingsStructure")
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

	stopTimeoutSecs, _ := valueObject.NewUnixTime(600)
	if serviceMap["stopTimeoutSecs"] != nil {
		stopTimeoutSecs, err = valueObject.NewUnixTime(
			serviceMap["stopTimeoutSecs"],
		)
		if err != nil {
			return installableService, err
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

	installTimeoutSecs, _ := valueObject.NewUnixTime(600)
	if serviceMap["installTimeoutSecs"] != nil {
		installTimeoutSecs, err = valueObject.NewUnixTime(
			serviceMap["installTimeoutSecs"],
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

	uninstallTimeoutSecs, _ := valueObject.NewUnixTime(600)
	if serviceMap["uninstallTimeoutSecs"] != nil {
		uninstallTimeoutSecs, err = valueObject.NewUnixTime(
			serviceMap["uninstallTimeoutSecs"],
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
			return installableService, errors.New("InvalidUninstallFilePathsStructure")
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

	preStartTimeoutSecs, _ := valueObject.NewUnixTime(600)
	if serviceMap["preStartTimeoutSecs"] != nil {
		preStartTimeoutSecs, err = valueObject.NewUnixTime(
			serviceMap["preStartTimeoutSecs"],
		)
		if err != nil {
			return installableService, err
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

	postStartTimeoutSecs, _ := valueObject.NewUnixTime(600)
	if serviceMap["postStartTimeoutSecs"] != nil {
		postStartTimeoutSecs, err = valueObject.NewUnixTime(
			serviceMap["postStartTimeoutSecs"],
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

	preStopTimeoutSecs, _ := valueObject.NewUnixTime(600)
	if serviceMap["preStopTimeoutSecs"] != nil {
		preStopTimeoutSecs, err = valueObject.NewUnixTime(
			serviceMap["preStopTimeoutSecs"],
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

	postStopTimeoutSecs, _ := valueObject.NewUnixTime(600)
	if serviceMap["postStopTimeoutSecs"] != nil {
		postStopTimeoutSecs, err = valueObject.NewUnixTime(
			serviceMap["postStopTimeoutSecs"],
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
		workingDirectory, err := valueObject.NewUnixFilePath(
			serviceMap["workingDirectory"],
		)
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
		manifestVersion, name, nature, serviceType, startCommand, description, versions,
		envs, portBindings, stopTimeoutSecs, stopCmdSteps, installTimeoutSecs,
		installCmdSteps, uninstallTimeoutSecs, uninstallCmdSteps, uninstallFilePaths,
		preStartTimeoutSecs, preStartCmdSteps, postStartTimeoutSecs, postStartCmdSteps,
		preStopTimeoutSecs, preStopCmdSteps, postStopTimeoutSecs, postStopCmdSteps,
		execUserPtr, workingDirectoryPtr, startupFilePtr, logOutputPathPtr,
		logErrorPathPtr, avatarUrlPtr, estimatedSizeBytesPtr,
	), nil
}

func (repo *ServicesQueryRepo) ReadInstallableItems(
	requestDto dto.ReadInstallableServicesItemsRequest,
) (installableItemsDto dto.ReadInstallableServicesItemsResponse, err error) {
	_, err = os.Stat(infraEnvs.InstallableServicesItemsDir)
	if err != nil {
		servicesCmdRepo := NewServicesCmdRepo(repo.persistentDbSvc)
		err = servicesCmdRepo.RefreshInstallableItems()
		if err != nil {
			return installableItemsDto, errors.New(
				"RefreshServiceInstallableItemsError: " + err.Error(),
			)
		}
	}

	rawInstallableFilesList, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "find " + infraEnvs.InstallableServicesItemsDir + " -type f " +
			"\\( -name '*.json' -o -name '*.yaml' -o -name '*.yml' \\) " +
			"-not -path '*/.*' -not -name '.*'",
		ShouldRunWithSubShell: true,
	})
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
			slog.Debug(err.Error(), slog.String("filePath", rawFilePath))
			continue
		}

		installableService, err := repo.installableServiceFactory(itemFilePath)
		if err != nil {
			slog.Debug(
				"CatalogMarketplaceItemFactoryError",
				slog.String("filePath", itemFilePath.String()),
				slog.String("err", err.Error()),
			)
			continue
		}

		installableServices = append(installableServices, installableService)
	}

	filteredInstallableServices := []entity.InstallableService{}
	for _, installableService := range installableServices {
		if requestDto.ServiceName != nil {
			if installableService.Name != *requestDto.ServiceName {
				continue
			}
		}

		if requestDto.ServiceNature != nil {
			if installableService.Nature != *requestDto.ServiceNature {
				continue
			}
		}

		if requestDto.ServiceType != nil && installableService.Type != *requestDto.ServiceType {
			if installableService.Type != *requestDto.ServiceType {
				continue
			}
		}

		filteredInstallableServices = append(
			filteredInstallableServices, installableService,
		)
	}

	if len(filteredInstallableServices) > int(requestDto.Pagination.ItemsPerPage) {
		filteredInstallableServices = filteredInstallableServices[:requestDto.Pagination.ItemsPerPage]
	}

	sortDirectionStr := "asc"
	if requestDto.Pagination.SortDirection != nil {
		sortDirectionStr = requestDto.Pagination.SortDirection.String()
	}

	if requestDto.Pagination.SortBy != nil {
		slices.SortStableFunc(filteredInstallableServices, func(a, b entity.InstallableService) int {
			firstElement := a
			secondElement := b
			if sortDirectionStr != "asc" {
				firstElement = b
				secondElement = a
			}

			switch requestDto.Pagination.SortBy.String() {
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

	itemsTotal := uint64(len(filteredInstallableServices))
	pagesTotal := uint32(itemsTotal / uint64(requestDto.Pagination.ItemsPerPage))

	paginationDto := requestDto.Pagination
	paginationDto.ItemsTotal = &itemsTotal
	paginationDto.PagesTotal = &pagesTotal

	return dto.ReadInstallableServicesItemsResponse{
		Pagination:          paginationDto,
		InstallableServices: filteredInstallableServices,
	}, nil
}

func (repo *ServicesQueryRepo) ReadFirstInstallableItem(
	requestDto dto.ReadInstallableServicesItemsRequest,
) (installableService entity.InstallableService, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.ReadInstallableItems(requestDto)
	if err != nil {
		return installableService, err
	}

	if len(responseDto.InstallableServices) == 0 {
		return installableService, errors.New("InstallableServiceItemNotFound")
	}

	return responseDto.InstallableServices[0], nil
}
