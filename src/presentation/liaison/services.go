package liaison

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
	tkInfra "github.com/goinfinite/tk/src/infra"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	liaisonHelper "github.com/goinfinite/os/src/presentation/liaison/helper"
)

type ServicesLiaison struct {
	persistentDbService   *internalDbInfra.PersistentDatabaseService
	servicesQueryRepo     *servicesInfra.ServicesQueryRepo
	servicesCmdRepo       *servicesInfra.ServicesCmdRepo
	mappingQueryRepo      *vhostInfra.MappingQueryRepo
	mappingCmdRepo        *vhostInfra.MappingCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewServicesLiaison(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *ServicesLiaison {
	return &ServicesLiaison{
		persistentDbService:   persistentDbService,
		servicesQueryRepo:     servicesInfra.NewServicesQueryRepo(persistentDbService),
		servicesCmdRepo:       servicesInfra.NewServicesCmdRepo(persistentDbService),
		mappingQueryRepo:      vhostInfra.NewMappingQueryRepo(persistentDbService),
		mappingCmdRepo:        vhostInfra.NewMappingCmdRepo(persistentDbService),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (liaison *ServicesLiaison) ReadInstalledItems(
	untrustedInput map[string]any,
) LiaisonOutput {
	var namePtr *valueObject.ServiceName
	if untrustedInput["name"] != nil {
		name, err := valueObject.NewServiceName(untrustedInput["name"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		namePtr = &name
	}

	var naturePtr *valueObject.ServiceNature
	if untrustedInput["nature"] != nil {
		nature, err := valueObject.NewServiceNature(untrustedInput["nature"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		naturePtr = &nature
	}

	var statusPtr *valueObject.ServiceStatus
	if untrustedInput["status"] != nil {
		status, err := valueObject.NewServiceStatus(untrustedInput["status"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		statusPtr = &status
	}

	var typePtr *valueObject.ServiceType
	if untrustedInput["type"] != nil {
		itemType, err := valueObject.NewServiceType(untrustedInput["type"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		typePtr = &itemType
	}

	shouldIncludeMetrics := false
	if untrustedInput["shouldIncludeMetrics"] != nil {
		var err error
		shouldIncludeMetrics, err = tkVoUtil.InterfaceToBool(untrustedInput["shouldIncludeMetrics"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
	}

	paginationDto := useCase.ServicesDefaultPagination
	if untrustedInput["pageNumber"] != nil {
		pageNumber, err := tkVoUtil.InterfaceToUint32(untrustedInput["pageNumber"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if untrustedInput["itemsPerPage"] != nil {
		itemsPerPage, err := tkVoUtil.InterfaceToUint16(untrustedInput["itemsPerPage"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if untrustedInput["sortBy"] != nil {
		sortBy, err := tkValueObject.NewPaginationSortBy(untrustedInput["sortBy"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if untrustedInput["sortDirection"] != nil {
		sortDirection, err := tkValueObject.NewPaginationSortDirection(
			untrustedInput["sortDirection"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if untrustedInput["lastSeenId"] != nil {
		lastSeenId, err := tkValueObject.NewPaginationLastSeenId(untrustedInput["lastSeenId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
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
		liaison.servicesQueryRepo, readDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, servicesList)
}

func (liaison *ServicesLiaison) ReadInstallableItems(
	untrustedInput map[string]any,
) LiaisonOutput {
	var namePtr *valueObject.ServiceName
	if untrustedInput["name"] != nil {
		name, err := valueObject.NewServiceName(untrustedInput["name"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		namePtr = &name
	}

	var naturePtr *valueObject.ServiceNature
	if untrustedInput["nature"] != nil {
		nature, err := valueObject.NewServiceNature(untrustedInput["nature"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		naturePtr = &nature
	}

	var typePtr *valueObject.ServiceType
	if untrustedInput["type"] != nil {
		itemType, err := valueObject.NewServiceType(untrustedInput["type"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		typePtr = &itemType
	}

	paginationDto := useCase.ServicesDefaultPagination
	if untrustedInput["pageNumber"] != nil {
		pageNumber, err := tkVoUtil.InterfaceToUint32(untrustedInput["pageNumber"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if untrustedInput["itemsPerPage"] != nil {
		itemsPerPage, err := tkVoUtil.InterfaceToUint16(untrustedInput["itemsPerPage"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if untrustedInput["sortBy"] != nil {
		sortBy, err := tkValueObject.NewPaginationSortBy(untrustedInput["sortBy"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if untrustedInput["sortDirection"] != nil {
		sortDirection, err := tkValueObject.NewPaginationSortDirection(
			untrustedInput["sortDirection"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if untrustedInput["lastSeenId"] != nil {
		lastSeenId, err := tkValueObject.NewPaginationLastSeenId(untrustedInput["lastSeenId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
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
		liaison.servicesQueryRepo, readDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, servicesList)
}

func (liaison *ServicesLiaison) CreateInstallable(
	untrustedInput map[string]any,
	shouldSchedule bool,
) LiaisonOutput {
	requiredParams := []string{"name"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	name, err := valueObject.NewServiceName(untrustedInput["name"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	var versionPtr *valueObject.ServiceVersion
	if untrustedInput["version"] != nil {
		version, err := valueObject.NewServiceVersion(untrustedInput["version"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		versionPtr = &version
	}

	var startupFilePtr *tkValueObject.UnixAbsoluteFilePath
	if untrustedInput["startupFile"] != nil {
		startupFile, err := tkValueObject.NewUnixAbsoluteFilePath(untrustedInput["startupFile"], false)
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		startupFilePtr = &startupFile
	}

	var workingDirPtr *tkValueObject.UnixAbsoluteFilePath
	if untrustedInput["workingDir"] != nil {
		workingDir, err := tkValueObject.NewUnixAbsoluteFilePath(untrustedInput["workingDir"], false)
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		workingDirPtr = &workingDir
	}

	envs := []valueObject.ServiceEnv{}
	if untrustedInput["envs"] != nil {
		var assertOk bool
		envs, assertOk = untrustedInput["envs"].([]valueObject.ServiceEnv)
		if !assertOk {
			return NewLiaisonOutput(UserError, "InvalidServiceEnvs")
		}
	}

	portBindings := []valueObject.PortBinding{}
	if untrustedInput["portBindings"] != nil {
		var assertOk bool
		portBindings, assertOk = untrustedInput["portBindings"].([]valueObject.PortBinding)
		if !assertOk {
			return NewLiaisonOutput(UserError, "InvalidPortBindings")
		}
	}

	var autoStartPtr *bool
	if untrustedInput["autoStart"] != nil {
		autoStart, err := tkVoUtil.InterfaceToBool(untrustedInput["autoStart"])
		if err != nil {
			return NewLiaisonOutput(UserError, "AutoStartMustBeBool")
		}
		autoStartPtr = &autoStart
	}

	var timeoutStartSecsPtr *uint
	if untrustedInput["timeoutStartSecs"] != nil {
		timeoutStartSecs, err := tkVoUtil.InterfaceToUint(untrustedInput["timeoutStartSecs"])
		if err != nil {
			return NewLiaisonOutput(UserError, "TimeoutStartSecsMustBeUint")
		}
		timeoutStartSecsPtr = &timeoutStartSecs
	}

	var autoRestartPtr *bool
	if untrustedInput["autoRestart"] != nil {
		autoRestart, err := tkVoUtil.InterfaceToBool(
			untrustedInput["autoRestart"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, "AutoRestartMustBeBool")
		}
		autoRestartPtr = &autoRestart
	}

	var maxStartRetriesPtr *uint
	if untrustedInput["maxStartRetries"] != nil {
		maxStartRetries, err := tkVoUtil.InterfaceToUint(untrustedInput["maxStartRetries"])
		if err != nil {
			return NewLiaisonOutput(UserError, "MaxStartRetriesMustBeUint")
		}
		maxStartRetriesPtr = &maxStartRetries
	}

	autoCreateMapping := true
	if untrustedInput["autoCreateMapping"] != nil {
		autoCreateMapping, err = tkVoUtil.InterfaceToBool(untrustedInput["autoCreateMapping"])
		if err != nil {
			return NewLiaisonOutput(UserError, "AutoCreateMappingMustBeBool")
		}
	}

	var mappingHostnamePtr *tkValueObject.Fqdn
	if untrustedInput["mappingHostname"] != nil {
		mappingHostname, err := tkValueObject.NewFqdn(untrustedInput["mappingHostname"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		mappingHostnamePtr = &mappingHostname
	}

	var mappingPathPtr *valueObject.MappingPath
	if untrustedInput["mappingPath"] != nil {
		mappingPath, err := valueObject.NewMappingPath(untrustedInput["mappingPath"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		mappingPathPtr = &mappingPath
	}

	var mappingUpgradeInsecureRequestsPtr *bool
	if untrustedInput["mappingUpgradeInsecureRequests"] != nil {
		mappingUpgradeInsecureRequests, err := tkVoUtil.InterfaceToBool(
			untrustedInput["mappingUpgradeInsecureRequests"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, "InvalidMappingUpgradeInsecureRequests")
		}
		mappingUpgradeInsecureRequestsPtr = &mappingUpgradeInsecureRequests
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
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
				escapedField := tkInfra.ShellEscape{}.Quote(env.String())
				installParams = append(installParams, "--envs", escapedField)
			}
		}

		if len(portBindings) > 0 {
			for _, portBinding := range portBindings {
				escapedField := tkInfra.ShellEscape{}.Quote(portBinding.String())
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

		scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(liaison.persistentDbService)
		taskName, _ := valueObject.NewScheduledTaskName("CreateInstallableService")
		taskCmd, _ := tkValueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("services")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSecs := uint16(1800)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSecs, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewLiaisonOutput(InfraError, err.Error())
		}

		return NewLiaisonOutput(Created, "CreateInstallableServiceScheduled")
	}

	createDto := dto.NewCreateInstallableService(
		name, envs, portBindings, versionPtr, startupFilePtr, workingDirPtr,
		autoStartPtr, timeoutStartSecsPtr, autoRestartPtr, maxStartRetriesPtr,
		&autoCreateMapping, mappingHostnamePtr, mappingPathPtr,
		mappingUpgradeInsecureRequestsPtr, operatorAccountId, operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(liaison.persistentDbService)

	err = useCase.CreateInstallableService(
		liaison.servicesQueryRepo, liaison.servicesCmdRepo, vhostQueryRepo,
		liaison.mappingCmdRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Created, "InstallableServiceCreated")
}

func (liaison *ServicesLiaison) CreateCustom(
	untrustedInput map[string]any,
) LiaisonOutput {
	requiredParams := []string{"name", "type", "startCmd"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	name, err := valueObject.NewServiceName(untrustedInput["name"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	svcType, err := valueObject.NewServiceType(untrustedInput["type"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	startCmd, err := tkValueObject.NewUnixCommand(untrustedInput["startCmd"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	var versionPtr *valueObject.ServiceVersion
	if untrustedInput["version"] != nil {
		version, err := valueObject.NewServiceVersion(untrustedInput["version"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		versionPtr = &version
	}

	var execUserPtr *tkValueObject.UnixUsername
	if untrustedInput["execUser"] != nil {
		execUser, err := tkValueObject.NewUnixUsername(untrustedInput["execUser"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		execUserPtr = &execUser
	}

	envs := []valueObject.ServiceEnv{}
	if untrustedInput["envs"] != nil {
		var assertOk bool
		envs, assertOk = untrustedInput["envs"].([]valueObject.ServiceEnv)
		if !assertOk {
			return NewLiaisonOutput(UserError, "InvalidServiceEnvs")
		}
	}

	portBindings := []valueObject.PortBinding{}
	if untrustedInput["portBindings"] != nil {
		var assertOk bool
		portBindings, assertOk = untrustedInput["portBindings"].([]valueObject.PortBinding)
		if !assertOk {
			return NewLiaisonOutput(UserError, "InvalidPortBindings")
		}
	}

	var autoStartPtr *bool
	if untrustedInput["autoStart"] != nil {
		autoStart, err := tkVoUtil.InterfaceToBool(untrustedInput["autoStart"])
		if err != nil {
			return NewLiaisonOutput(UserError, "AutoStartMustBeBool")
		}
		autoStartPtr = &autoStart
	}

	var timeoutStartSecsPtr *uint
	if untrustedInput["timeoutStartSecs"] != nil {
		timeoutStartSecs, err := tkVoUtil.InterfaceToUint(untrustedInput["timeoutStartSecs"])
		if err != nil {
			return NewLiaisonOutput(UserError, "TimeoutStartSecsMustBeUint")
		}
		timeoutStartSecsPtr = &timeoutStartSecs
	}

	var autoRestartPtr *bool
	if untrustedInput["autoRestart"] != nil {
		autoRestart, err := tkVoUtil.InterfaceToBool(
			untrustedInput["autoRestart"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, "AutoRestartMustBeBool")
		}
		autoRestartPtr = &autoRestart
	}

	var maxStartRetriesPtr *uint
	if untrustedInput["maxStartRetries"] != nil {
		maxStartRetries, err := tkVoUtil.InterfaceToUint(untrustedInput["maxStartRetries"])
		if err != nil {
			return NewLiaisonOutput(UserError, "MaxStartRetriesMustBeUint")
		}
		maxStartRetriesPtr = &maxStartRetries
	}

	var logOutputPathPtr *tkValueObject.UnixAbsoluteFilePath
	if untrustedInput["logOutputPath"] != nil {
		logOutputPath, err := tkValueObject.NewUnixAbsoluteFilePath(untrustedInput["logOutputPath"], false)
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		logOutputPathPtr = &logOutputPath
	}

	var logErrorPathPtr *tkValueObject.UnixAbsoluteFilePath
	if untrustedInput["logErrorPath"] != nil {
		logErrorPath, err := tkValueObject.NewUnixAbsoluteFilePath(untrustedInput["logErrorPath"], false)
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		logErrorPathPtr = &logErrorPath
	}

	var avatarUrlPtr *tkValueObject.Url
	if untrustedInput["avatarUrl"] != nil {
		avatarUrl, err := tkValueObject.NewUrl(untrustedInput["avatarUrl"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		avatarUrlPtr = &avatarUrl
	}

	autoCreateMapping := true
	if untrustedInput["autoCreateMapping"] != nil {
		autoCreateMapping, err = tkVoUtil.InterfaceToBool(untrustedInput["autoCreateMapping"])
		if err != nil {
			return NewLiaisonOutput(UserError, "AutoCreateMappingMustBeBool")
		}
	}

	var mappingHostnamePtr *tkValueObject.Fqdn
	if untrustedInput["mappingHostname"] != nil {
		mappingHostname, err := tkValueObject.NewFqdn(untrustedInput["mappingHostname"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		mappingHostnamePtr = &mappingHostname
	}

	var mappingPathPtr *valueObject.MappingPath
	if untrustedInput["mappingPath"] != nil {
		mappingPath, err := valueObject.NewMappingPath(untrustedInput["mappingPath"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		mappingPathPtr = &mappingPath
	}

	var mappingUpgradeInsecureRequestsPtr *bool
	if untrustedInput["mappingUpgradeInsecureRequests"] != nil {
		mappingUpgradeInsecureRequests, err := tkVoUtil.InterfaceToBool(
			untrustedInput["mappingUpgradeInsecureRequests"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, "InvalidMappingUpgradeInsecureRequests")
		}
		mappingUpgradeInsecureRequestsPtr = &mappingUpgradeInsecureRequests
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	createCustomDto := dto.NewCreateCustomService(
		name, svcType, startCmd, envs, portBindings, nil, nil, nil, nil, nil,
		versionPtr, execUserPtr, nil, autoStartPtr, autoRestartPtr,
		timeoutStartSecsPtr, maxStartRetriesPtr, logOutputPathPtr, logErrorPathPtr,
		avatarUrlPtr, &autoCreateMapping, mappingHostnamePtr, mappingPathPtr,
		mappingUpgradeInsecureRequestsPtr, operatorAccountId, operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(liaison.persistentDbService)

	err = useCase.CreateCustomService(
		liaison.servicesQueryRepo, liaison.servicesCmdRepo, vhostQueryRepo,
		liaison.mappingCmdRepo, liaison.activityRecordCmdRepo, createCustomDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Created, "CustomServiceCreated")
}

func (liaison *ServicesLiaison) Update(untrustedInput map[string]any) LiaisonOutput {
	requiredParams := []string{"name"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	name, err := valueObject.NewServiceName(untrustedInput["name"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	var typePtr *valueObject.ServiceType
	if untrustedInput["type"] != nil {
		svcType, err := valueObject.NewServiceType(untrustedInput["type"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		typePtr = &svcType
	}

	var startCmdPtr *tkValueObject.UnixCommand
	if untrustedInput["startCmd"] != nil {
		startCmd, err := tkValueObject.NewUnixCommand(untrustedInput["startCmd"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		startCmdPtr = &startCmd
	}

	var statusPtr *valueObject.ServiceStatus
	if untrustedInput["status"] != nil {
		status, err := valueObject.NewServiceStatus(untrustedInput["status"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		statusPtr = &status
	}

	var versionPtr *valueObject.ServiceVersion
	if untrustedInput["version"] != nil {
		version, err := valueObject.NewServiceVersion(untrustedInput["version"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		versionPtr = &version
	}

	envs := []valueObject.ServiceEnv{}
	if untrustedInput["envs"] != nil {
		rawEnvs, assertOk := untrustedInput["envs"].([]string)
		if !assertOk {
			return NewLiaisonOutput(UserError, "EnvsMustBeStringArray")
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
	if _, exists := untrustedInput["portBindings"]; exists {
		rawPortBindings, assertOk := untrustedInput["portBindings"].([]string)
		if !assertOk {
			return NewLiaisonOutput(UserError, "PortBindingsMustBeStringArray")
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

	var startupFilePtr *tkValueObject.UnixAbsoluteFilePath
	if untrustedInput["startupFile"] != nil {
		startupFile, err := tkValueObject.NewUnixAbsoluteFilePath(untrustedInput["startupFile"], false)
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		startupFilePtr = &startupFile
	}

	var autoStartPtr *bool
	if untrustedInput["autoStart"] != nil {
		autoStart, err := tkVoUtil.InterfaceToBool(untrustedInput["autoStart"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		autoStartPtr = &autoStart
	}

	var autoRestartPtr *bool
	if untrustedInput["autoRestart"] != nil {
		autoRestart, err := tkVoUtil.InterfaceToBool(untrustedInput["autoRestart"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		autoRestartPtr = &autoRestart
	}

	var timeoutStartSecsPtr *uint
	if untrustedInput["timeoutStartSecs"] != nil {
		timeoutStartSecs, err := tkVoUtil.InterfaceToUint(untrustedInput["timeoutStartSecs"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		timeoutStartSecsPtr = &timeoutStartSecs
	}

	var maxStartRetriesPtr *uint
	if untrustedInput["maxStartRetries"] != nil {
		maxStartRetries, err := tkVoUtil.InterfaceToUint(untrustedInput["maxStartRetries"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		maxStartRetriesPtr = &maxStartRetries
	}

	var logOutputPathPtr *tkValueObject.UnixAbsoluteFilePath
	if untrustedInput["logOutputPath"] != nil {
		logOutputPath, err := tkValueObject.NewUnixAbsoluteFilePath(untrustedInput["logOutputPath"], false)
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		logOutputPathPtr = &logOutputPath
	}

	var logErrorPathPtr *tkValueObject.UnixAbsoluteFilePath
	if untrustedInput["logErrorPath"] != nil {
		logErrorPath, err := tkValueObject.NewUnixAbsoluteFilePath(untrustedInput["logErrorPath"], false)
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		logErrorPathPtr = &logErrorPath
	}

	var avatarUrlPtr *tkValueObject.Url
	if untrustedInput["avatarUrl"] != nil {
		avatarUrl, err := tkValueObject.NewUrl(untrustedInput["avatarUrl"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		avatarUrlPtr = &avatarUrl
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdateService(
		name, typePtr, versionPtr, statusPtr, startCmdPtr, envs, portBindings, nil,
		nil, nil, nil, nil, nil, nil, startupFilePtr, autoStartPtr, autoRestartPtr,
		timeoutStartSecsPtr, maxStartRetriesPtr, logOutputPathPtr, logErrorPathPtr,
		avatarUrlPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.UpdateService(
		liaison.servicesQueryRepo, liaison.servicesCmdRepo, liaison.mappingQueryRepo,
		liaison.mappingCmdRepo, liaison.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "ServiceUpdated")
}

func (liaison *ServicesLiaison) Delete(untrustedInput map[string]any) LiaisonOutput {
	requiredParams := []string{"name"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	name, err := valueObject.NewServiceName(untrustedInput["name"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteService(name, operatorAccountId, operatorIpAddress)

	err = useCase.DeleteService(
		liaison.servicesQueryRepo, liaison.servicesCmdRepo, liaison.mappingQueryRepo,
		liaison.mappingCmdRepo, liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "ServiceDeleted")
}
