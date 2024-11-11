package service

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	mappingInfra "github.com/goinfinite/os/src/infra/vhost/mapping"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
)

type ServicesService struct {
	persistentDbService   *internalDbInfra.PersistentDatabaseService
	servicesQueryRepo     *servicesInfra.ServicesQueryRepo
	servicesCmdRepo       *servicesInfra.ServicesCmdRepo
	mappingQueryRepo      *mappingInfra.MappingQueryRepo
	mappingCmdRepo        *mappingInfra.MappingCmdRepo
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
		mappingQueryRepo:      mappingInfra.NewMappingQueryRepo(persistentDbService),
		mappingCmdRepo:        mappingInfra.NewMappingCmdRepo(persistentDbService),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *ServicesService) Read() ServiceOutput {
	servicesList, err := useCase.ReadServicesWithMetrics(service.servicesQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, servicesList)
}

func (service *ServicesService) ReadInstallables() ServiceOutput {
	servicesList, err := useCase.ReadInstallableServices(service.servicesQueryRepo)
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
	if input["portBindings"] != nil {
		rawPortBindings, assertOk := input["portBindings"].([]string)
		if !assertOk {
			return NewServiceOutput(UserError, "PortBindingsMustBeStringArray")
		}

		for _, rawPortBinding := range rawPortBindings {
			portBinding, err := valueObject.NewPortBinding(rawPortBinding)
			if err != nil {
				slog.Debug(err.Error(), slog.String("portBinding", rawPortBinding))
				continue
			}
			portBindings = append(portBindings, portBinding)
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
				escapedField := shellescape.Quote(env.String())
				installParams = append(installParams, "--envs", escapedField)
			}
		}

		if len(portBindings) > 0 {
			for _, portBinding := range portBindings {
				escapedField := shellescape.Quote(portBinding.String())
				installParams = append(installParams, "--port-bindings", escapedField)
			}
		}

		if versionPtr != nil {
			installParams = append(installParams, "--version", versionPtr.String())
		}

		if startupFilePtr != nil {
			installParams = append(installParams, "--startup-file", startupFilePtr.String())
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

		cliCmd += " " + strings.Join(installParams, " ")

		scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(service.persistentDbService)
		taskName, _ := valueObject.NewScheduledTaskName("CreateInstallableService")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("services")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSeconds := uint16(600)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSeconds, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "CreateInstallableServiceScheduled")
	}

	createDto := dto.NewCreateInstallableService(
		name, envs, portBindings, versionPtr, startupFilePtr, autoStartPtr,
		timeoutStartSecsPtr, autoRestartPtr, maxStartRetriesPtr, &autoCreateMapping,
		operatorAccountId, operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbService)

	err = useCase.CreateInstallableService(
		service.servicesQueryRepo, service.servicesCmdRepo, service.mappingQueryRepo,
		service.mappingCmdRepo, vhostQueryRepo, service.activityRecordCmdRepo,
		createDto,
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

	portBindings := []valueObject.PortBinding{}
	if _, exists := input["portBindings"]; exists {
		rawPortBindings, assertOk := input["portBindings"].([]string)
		if !assertOk {
			return NewServiceOutput(UserError, "PortBindingsMustBeStringArray")
		}

		for _, rawPortBinding := range rawPortBindings {
			portBinding, err := valueObject.NewPortBinding(rawPortBinding)
			if err != nil {
				slog.Debug(err.Error(), slog.String("portBinding", rawPortBinding))
				continue
			}
			portBindings = append(portBindings, portBinding)
		}
	}

	autoCreateMapping := true
	if input["autoCreateMapping"] != nil {
		autoCreateMapping, err = voHelper.InterfaceToBool(input["autoCreateMapping"])
		if err != nil {
			return NewServiceOutput(UserError, "AutoCreateMappingMustBeBool")
		}
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
		name, svcType, startCmd, []valueObject.ServiceEnv{}, portBindings,
		nil, nil, nil, nil, nil, versionPtr, nil, nil, nil, nil, nil, nil, nil, nil,
		&autoCreateMapping, operatorAccountId, operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbService)

	err = useCase.CreateCustomService(
		service.servicesQueryRepo, service.servicesCmdRepo, service.mappingQueryRepo,
		service.mappingCmdRepo, vhostQueryRepo, service.activityRecordCmdRepo,
		createCustomDto,
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

	portBindings := []valueObject.PortBinding{}
	if _, exists := input["portBindings"]; exists {
		rawPortBindings, assertOk := input["portBindings"].([]string)
		if !assertOk {
			return NewServiceOutput(UserError, "PortBindingsMustBeStringArray")
		}

		for _, rawPortBinding := range rawPortBindings {
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

	dto := dto.NewUpdateService(
		name, typePtr, versionPtr, statusPtr, startCmdPtr, nil, portBindings, nil,
		nil, nil, nil, nil, nil, nil, startupFilePtr, nil, nil, nil, nil, nil, nil,
	)

	err = useCase.UpdateService(
		service.servicesQueryRepo, service.servicesCmdRepo, service.mappingQueryRepo,
		service.mappingCmdRepo, dto,
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

	err = useCase.DeleteService(
		service.servicesQueryRepo, service.servicesCmdRepo,
		service.mappingCmdRepo, name,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ServiceDeleted")
}
