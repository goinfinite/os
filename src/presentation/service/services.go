package service

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
)

type ServicesService struct {
	persistentDbService   *internalDbInfra.PersistentDatabaseService
	servicesQueryRepo     *servicesInfra.ServicesQueryRepo
	servicesCmdRepo       *servicesInfra.ServicesCmdRepo
	mappingQueryRepo      *vhostInfra.MappingQueryRepo
	mappingCmdRepo        *vhostInfra.MappingCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewServicesService(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *ServicesService {
	return &ServicesService{
		persistentDbService:   persistentDbService,
		servicesQueryRepo:     servicesInfra.NewServicesQueryRepo(persistentDbService),
		servicesCmdRepo:       servicesInfra.NewServicesCmdRepo(persistentDbService),
		mappingQueryRepo:      vhostInfra.NewMappingQueryRepo(persistentDbService),
		mappingCmdRepo:        vhostInfra.NewMappingCmdRepo(persistentDbService),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *ServicesService) ReadInstalledItems(
	input map[string]interface{},
) ServiceOutput {
	var namePtr *valueObject.ServiceName
	if input["name"] != nil {
		name, err := valueObject.NewServiceName(input["name"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		namePtr = &name
	}

	var naturePtr *valueObject.ServiceNature
	if input["nature"] != nil {
		nature, err := valueObject.NewServiceNature(input["nature"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		naturePtr = &nature
	}

	var statusPtr *valueObject.ServiceStatus
	if input["status"] != nil {
		status, err := valueObject.NewServiceStatus(input["status"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		statusPtr = &status
	}

	var typePtr *valueObject.ServiceType
	if input["type"] != nil {
		itemType, err := valueObject.NewServiceType(input["type"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		typePtr = &itemType
	}

	shouldIncludeMetrics := false
	if input["shouldIncludeMetrics"] != nil {
		var err error
		shouldIncludeMetrics, err = voHelper.InterfaceToBool(input["shouldIncludeMetrics"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
	}

	paginationDto := useCase.ServicesDefaultPagination
	if input["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(input["pageNumber"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if input["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(input["itemsPerPage"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if input["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(input["sortBy"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if input["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(
			input["sortDirection"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if input["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(input["lastSeenId"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadInstalledServicesItemsRequest{
		Pagination:           paginationDto,
		ServiceName:          namePtr,
		ServiceNature:        naturePtr,
		ServiceType:          typePtr,
		ServiceStatus:        statusPtr,
		ShouldIncludeMetrics: &shouldIncludeMetrics,
	}

	servicesList, err := useCase.ReadInstalledServices(
		service.servicesQueryRepo, readDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, servicesList)
}

func (service *ServicesService) ReadInstallableItems(
	input map[string]interface{},
) ServiceOutput {
	var namePtr *valueObject.ServiceName
	if input["name"] != nil {
		name, err := valueObject.NewServiceName(input["name"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		namePtr = &name
	}

	var naturePtr *valueObject.ServiceNature
	if input["nature"] != nil {
		nature, err := valueObject.NewServiceNature(input["nature"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		naturePtr = &nature
	}

	var typePtr *valueObject.ServiceType
	if input["type"] != nil {
		itemType, err := valueObject.NewServiceType(input["type"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		typePtr = &itemType
	}

	paginationDto := useCase.ServicesDefaultPagination
	if input["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(input["pageNumber"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if input["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(input["itemsPerPage"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if input["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(input["sortBy"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if input["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(
			input["sortDirection"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if input["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(input["lastSeenId"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadInstallableServicesItemsRequest{
		Pagination:    paginationDto,
		ServiceName:   namePtr,
		ServiceNature: naturePtr,
		ServiceType:   typePtr,
	}

	servicesList, err := useCase.ReadInstallableServices(
		service.servicesQueryRepo, readDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, servicesList)
}

func (service *ServicesService) CreateInstallable(
	input map[string]interface{},
	shouldSchedule bool,
) ServiceOutput {
	requiredParams := []string{"name"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	name, err := valueObject.NewServiceName(input["name"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var versionPtr *valueObject.ServiceVersion
	if input["version"] != nil {
		version, err := valueObject.NewServiceVersion(input["version"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		versionPtr = &version
	}

	var startupFilePtr *valueObject.UnixFilePath
	if input["startupFile"] != nil {
		startupFile, err := valueObject.NewUnixFilePath(input["startupFile"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		startupFilePtr = &startupFile
	}

	var workingDirPtr *valueObject.UnixFilePath
	if input["workingDir"] != nil {
		workingDir, err := valueObject.NewUnixFilePath(input["workingDir"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		workingDirPtr = &workingDir
	}

	envs := []valueObject.ServiceEnv{}
	if input["envs"] != nil {
		var assertOk bool
		envs, assertOk = input["envs"].([]valueObject.ServiceEnv)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidServiceEnvs")
		}
	}

	portBindings := []valueObject.PortBinding{}
	if input["portBindings"] != nil {
		var assertOk bool
		portBindings, assertOk = input["portBindings"].([]valueObject.PortBinding)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidPortBindings")
		}
	}

	var autoStartPtr *bool
	if input["autoStart"] != nil {
		autoStart, err := voHelper.InterfaceToBool(input["autoStart"])
		if err != nil {
			return NewServiceOutput(UserError, "AutoStartMustBeBool")
		}
		autoStartPtr = &autoStart
	}

	var timeoutStartSecsPtr *uint
	if input["timeoutStartSecs"] != nil {
		timeoutStartSecs, err := voHelper.InterfaceToUint(input["timeoutStartSecs"])
		if err != nil {
			return NewServiceOutput(UserError, "TimeoutStartSecsMustBeUint")
		}
		timeoutStartSecsPtr = &timeoutStartSecs
	}

	var autoRestartPtr *bool
	if input["autoRestart"] != nil {
		autoRestart, err := voHelper.InterfaceToBool(
			input["autoRestart"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "AutoRestartMustBeBool")
		}
		autoRestartPtr = &autoRestart
	}

	var maxStartRetriesPtr *uint
	if input["maxStartRetries"] != nil {
		maxStartRetries, err := voHelper.InterfaceToUint(input["maxStartRetries"])
		if err != nil {
			return NewServiceOutput(UserError, "MaxStartRetriesMustBeUint")
		}
		maxStartRetriesPtr = &maxStartRetries
	}

	autoCreateMapping := true
	if input["autoCreateMapping"] != nil {
		autoCreateMapping, err = voHelper.InterfaceToBool(input["autoCreateMapping"])
		if err != nil {
			return NewServiceOutput(UserError, "AutoCreateMappingMustBeBool")
		}
	}

	var mappingHostnamePtr *valueObject.Fqdn
	if input["mappingHostname"] != nil {
		mappingHostname, err := valueObject.NewFqdn(input["mappingHostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		mappingHostnamePtr = &mappingHostname
	}

	var mappingPathPtr *valueObject.MappingPath
	if input["mappingPath"] != nil {
		mappingPath, err := valueObject.NewMappingPath(input["mappingPath"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		mappingPathPtr = &mappingPath
	}

	var mappingUpgradeInsecureRequestsPtr *bool
	if input["mappingUpgradeInsecureRequests"] != nil {
		mappingUpgradeInsecureRequests, err := voHelper.InterfaceToBool(
			input["mappingUpgradeInsecureRequests"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMappingUpgradeInsecureRequests")
		}
		mappingUpgradeInsecureRequestsPtr = &mappingUpgradeInsecureRequests
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteOsBinary + " services create-installable"
		installParams := []string{
			"--name", name.String(),
			"--auto-create-mapping", strconv.FormatBool(autoCreateMapping),
		}

		if len(envs) > 0 {
			for _, env := range envs {
				escapedField := infraHelper.ShellEscape{}.Quote(env.String())
				installParams = append(installParams, "--envs", escapedField)
			}
		}

		if len(portBindings) > 0 {
			for _, portBinding := range portBindings {
				escapedField := infraHelper.ShellEscape{}.Quote(portBinding.String())
				installParams = append(installParams, "--port-bindings", escapedField)
			}
		}

		if versionPtr != nil {
			installParams = append(installParams, "--version", versionPtr.String())
		}

		if startupFilePtr != nil {
			installParams = append(installParams, "--startup-file", startupFilePtr.String())
		}

		if workingDirPtr != nil {
			installParams = append(installParams, "--working-dir", workingDirPtr.String())
		}

		if autoStartPtr != nil {
			autoStartStr := strconv.FormatBool(*autoStartPtr)
			installParams = append(installParams, "--auto-start", autoStartStr)
		}

		if timeoutStartSecsPtr != nil {
			timeoutStartSecsStr := strconv.FormatUint(uint64(*timeoutStartSecsPtr), 10)
			installParams = append(installParams, "--timeout-start-secs", timeoutStartSecsStr)
		}

		if autoRestartPtr != nil {
			autoRestartStr := strconv.FormatBool(*autoRestartPtr)
			installParams = append(installParams, "--auto-restart", autoRestartStr)
		}

		if maxStartRetriesPtr != nil {
			maxStartRetriesStr := strconv.FormatUint(uint64(*maxStartRetriesPtr), 10)
			installParams = append(installParams, "--max-start-retries", maxStartRetriesStr)
		}

		if mappingHostnamePtr != nil {
			installParams = append(installParams, "--mapping-hostname", mappingHostnamePtr.String())
		}

		if mappingPathPtr != nil {
			installParams = append(installParams, "--mapping-path", mappingPathPtr.String())
		}

		if mappingUpgradeInsecureRequestsPtr != nil {
			installParams = append(
				installParams,
				"--mapping-upgrade-insecure-requests",
				strconv.FormatBool(*mappingUpgradeInsecureRequestsPtr),
			)
		}

		cliCmd += " " + strings.Join(installParams, " ")

		scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(service.persistentDbService)
		taskName, _ := valueObject.NewScheduledTaskName("CreateInstallableService")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("services")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSecs := uint16(1800)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSecs, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "CreateInstallableServiceScheduled")
	}

	createDto := dto.NewCreateInstallableService(
		name, envs, portBindings, versionPtr, startupFilePtr, workingDirPtr,
		autoStartPtr, timeoutStartSecsPtr, autoRestartPtr, maxStartRetriesPtr,
		&autoCreateMapping, mappingHostnamePtr, mappingPathPtr,
		mappingUpgradeInsecureRequestsPtr, operatorAccountId, operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbService)

	err = useCase.CreateInstallableService(
		service.servicesQueryRepo, service.servicesCmdRepo, vhostQueryRepo,
		service.mappingCmdRepo, service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "InstallableServiceCreated")
}

func (service *ServicesService) CreateCustom(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"name", "type", "startCmd"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	name, err := valueObject.NewServiceName(input["name"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	svcType, err := valueObject.NewServiceType(input["type"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	startCmd, err := valueObject.NewUnixCommand(input["startCmd"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var versionPtr *valueObject.ServiceVersion
	if input["version"] != nil {
		version, err := valueObject.NewServiceVersion(input["version"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		versionPtr = &version
	}

	var execUserPtr *valueObject.UnixUsername
	if input["execUser"] != nil {
		execUser, err := valueObject.NewUnixUsername(input["execUser"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		execUserPtr = &execUser
	}

	envs := []valueObject.ServiceEnv{}
	if input["envs"] != nil {
		var assertOk bool
		envs, assertOk = input["envs"].([]valueObject.ServiceEnv)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidServiceEnvs")
		}
	}

	portBindings := []valueObject.PortBinding{}
	if input["portBindings"] != nil {
		var assertOk bool
		portBindings, assertOk = input["portBindings"].([]valueObject.PortBinding)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidPortBindings")
		}
	}

	var autoStartPtr *bool
	if input["autoStart"] != nil {
		autoStart, err := voHelper.InterfaceToBool(input["autoStart"])
		if err != nil {
			return NewServiceOutput(UserError, "AutoStartMustBeBool")
		}
		autoStartPtr = &autoStart
	}

	var timeoutStartSecsPtr *uint
	if input["timeoutStartSecs"] != nil {
		timeoutStartSecs, err := voHelper.InterfaceToUint(input["timeoutStartSecs"])
		if err != nil {
			return NewServiceOutput(UserError, "TimeoutStartSecsMustBeUint")
		}
		timeoutStartSecsPtr = &timeoutStartSecs
	}

	var autoRestartPtr *bool
	if input["autoRestart"] != nil {
		autoRestart, err := voHelper.InterfaceToBool(
			input["autoRestart"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "AutoRestartMustBeBool")
		}
		autoRestartPtr = &autoRestart
	}

	var maxStartRetriesPtr *uint
	if input["maxStartRetries"] != nil {
		maxStartRetries, err := voHelper.InterfaceToUint(input["maxStartRetries"])
		if err != nil {
			return NewServiceOutput(UserError, "MaxStartRetriesMustBeUint")
		}
		maxStartRetriesPtr = &maxStartRetries
	}

	var logOutputPathPtr *valueObject.UnixFilePath
	if input["logOutputPath"] != nil {
		logOutputPath, err := valueObject.NewUnixFilePath(input["logOutputPath"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		logOutputPathPtr = &logOutputPath
	}

	var logErrorPathPtr *valueObject.UnixFilePath
	if input["logErrorPath"] != nil {
		logErrorPath, err := valueObject.NewUnixFilePath(input["logErrorPath"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		logErrorPathPtr = &logErrorPath
	}

	var avatarUrlPtr *valueObject.Url
	if input["avatarUrl"] != nil {
		avatarUrl, err := valueObject.NewUrl(input["avatarUrl"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		avatarUrlPtr = &avatarUrl
	}

	autoCreateMapping := true
	if input["autoCreateMapping"] != nil {
		autoCreateMapping, err = voHelper.InterfaceToBool(input["autoCreateMapping"])
		if err != nil {
			return NewServiceOutput(UserError, "AutoCreateMappingMustBeBool")
		}
	}

	var mappingHostnamePtr *valueObject.Fqdn
	if input["mappingHostname"] != nil {
		mappingHostname, err := valueObject.NewFqdn(input["mappingHostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		mappingHostnamePtr = &mappingHostname
	}

	var mappingPathPtr *valueObject.MappingPath
	if input["mappingPath"] != nil {
		mappingPath, err := valueObject.NewMappingPath(input["mappingPath"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		mappingPathPtr = &mappingPath
	}

	var mappingUpgradeInsecureRequestsPtr *bool
	if input["mappingUpgradeInsecureRequests"] != nil {
		mappingUpgradeInsecureRequests, err := voHelper.InterfaceToBool(
			input["mappingUpgradeInsecureRequests"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidMappingUpgradeInsecureRequests")
		}
		mappingUpgradeInsecureRequestsPtr = &mappingUpgradeInsecureRequests
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createCustomDto := dto.NewCreateCustomService(
		name, svcType, startCmd, envs, portBindings, nil, nil, nil, nil, nil,
		versionPtr, execUserPtr, nil, autoStartPtr, autoRestartPtr,
		timeoutStartSecsPtr, maxStartRetriesPtr, logOutputPathPtr, logErrorPathPtr,
		avatarUrlPtr, &autoCreateMapping, mappingHostnamePtr, mappingPathPtr,
		mappingUpgradeInsecureRequestsPtr, operatorAccountId, operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbService)

	err = useCase.CreateCustomService(
		service.servicesQueryRepo, service.servicesCmdRepo, vhostQueryRepo,
		service.mappingCmdRepo, service.activityRecordCmdRepo, createCustomDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "CustomServiceCreated")
}

func (service *ServicesService) Update(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"name"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	name, err := valueObject.NewServiceName(input["name"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var typePtr *valueObject.ServiceType
	if input["type"] != nil {
		svcType, err := valueObject.NewServiceType(input["type"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		typePtr = &svcType
	}

	var startCmdPtr *valueObject.UnixCommand
	if input["startCmd"] != nil {
		startCmd, err := valueObject.NewUnixCommand(input["startCmd"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		startCmdPtr = &startCmd
	}

	var statusPtr *valueObject.ServiceStatus
	if input["status"] != nil {
		status, err := valueObject.NewServiceStatus(input["status"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		statusPtr = &status
	}

	var versionPtr *valueObject.ServiceVersion
	if input["version"] != nil {
		version, err := valueObject.NewServiceVersion(input["version"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		versionPtr = &version
	}

	envs := []valueObject.ServiceEnv{}
	if input["envs"] != nil {
		rawEnvs, assertOk := input["envs"].([]string)
		if !assertOk {
			return NewServiceOutput(UserError, "EnvsMustBeStringArray")
		}

		for _, rawEnv := range rawEnvs {
			env, err := valueObject.NewServiceEnv(rawEnv)
			if err != nil {
				slog.Debug(err.Error(), slog.String("env", rawEnv))
				continue
			}
			envs = append(envs, env)
		}
	}

	portBindings := []valueObject.PortBinding{}
	if _, exists := input["portBindings"]; exists {
		rawPortBindings, assertOk := input["portBindings"].([]string)
		if !assertOk {
			return NewServiceOutput(UserError, "PortBindingsMustBeStringArray")
		}

		for _, rawPortBinding := range rawPortBindings {
			if len(rawPortBinding) == 0 {
				continue
			}

			portBinding, err := valueObject.NewPortBinding(rawPortBinding)
			if err != nil {
				slog.Debug(err.Error(), slog.String("portBinding", rawPortBinding))
				continue
			}
			portBindings = append(portBindings, portBinding)
		}
	}

	var startupFilePtr *valueObject.UnixFilePath
	if input["startupFile"] != nil {
		startupFile, err := valueObject.NewUnixFilePath(input["startupFile"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		startupFilePtr = &startupFile
	}

	var autoStartPtr *bool
	if input["autoStart"] != nil {
		autoStart, err := voHelper.InterfaceToBool(input["autoStart"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		autoStartPtr = &autoStart
	}

	var autoRestartPtr *bool
	if input["autoRestart"] != nil {
		autoRestart, err := voHelper.InterfaceToBool(input["autoRestart"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		autoRestartPtr = &autoRestart
	}

	var timeoutStartSecsPtr *uint
	if input["timeoutStartSecs"] != nil {
		timeoutStartSecs, err := voHelper.InterfaceToUint(input["timeoutStartSecs"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		timeoutStartSecsPtr = &timeoutStartSecs
	}

	var maxStartRetriesPtr *uint
	if input["maxStartRetries"] != nil {
		maxStartRetries, err := voHelper.InterfaceToUint(input["maxStartRetries"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		maxStartRetriesPtr = &maxStartRetries
	}

	var logOutputPathPtr *valueObject.UnixFilePath
	if input["logOutputPath"] != nil {
		logOutputPath, err := valueObject.NewUnixFilePath(input["logOutputPath"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		logOutputPathPtr = &logOutputPath
	}

	var logErrorPathPtr *valueObject.UnixFilePath
	if input["logErrorPath"] != nil {
		logErrorPath, err := valueObject.NewUnixFilePath(input["logErrorPath"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		logErrorPathPtr = &logErrorPath
	}

	var avatarUrlPtr *valueObject.Url
	if input["avatarUrl"] != nil {
		avatarUrl, err := valueObject.NewUrl(input["avatarUrl"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		avatarUrlPtr = &avatarUrl
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdateService(
		name, typePtr, versionPtr, statusPtr, startCmdPtr, envs, portBindings, nil,
		nil, nil, nil, nil, nil, nil, startupFilePtr, autoStartPtr, autoRestartPtr,
		timeoutStartSecsPtr, maxStartRetriesPtr, logOutputPathPtr, logErrorPathPtr,
		avatarUrlPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.UpdateService(
		service.servicesQueryRepo, service.servicesCmdRepo, service.mappingQueryRepo,
		service.mappingCmdRepo, service.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ServiceUpdated")
}

func (service *ServicesService) Delete(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"name"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	name, err := valueObject.NewServiceName(input["name"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteService(name, operatorAccountId, operatorIpAddress)

	err = useCase.DeleteService(
		service.servicesQueryRepo, service.servicesCmdRepo, service.mappingQueryRepo,
		service.mappingCmdRepo, service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ServiceDeleted")
}
