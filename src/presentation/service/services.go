package service

import (
	"log"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraEnvs "github.com/speedianet/os/src/infra/envs"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/speedianet/os/src/infra/scheduledTask"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
)

type ServicesService struct {
	persistentDbService *internalDbInfra.PersistentDatabaseService
}

func NewServicesService(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
) *ServicesService {
	return &ServicesService{
		persistentDbService: persistentDbService,
	}
}

func (service *ServicesService) Read() ServiceOutput {
	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(service.persistentDbService)
	servicesList, err := useCase.ReadServicesWithMetrics(servicesQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, servicesList)
}

func (service *ServicesService) ReadInstallables() ServiceOutput {
	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(service.persistentDbService)
	servicesList, err := useCase.ReadInstallableServices(servicesQueryRepo)
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
			return NewServiceOutput(UserError, "InvalidEnvs")
		}

		for _, rawEnv := range rawEnvs {
			env, err := valueObject.NewServiceEnv(rawEnv)
			if err != nil {
				continue
			}
			envs = append(envs, env)
		}
	}

	portBindings := []valueObject.PortBinding{}
	if input["portBindings"] != nil {
		rawPortBindings, assertOk := input["portBindings"].([]string)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidPortBindings")
		}

		for _, rawPortBinding := range rawPortBindings {
			portBinding, err := valueObject.NewPortBinding(rawPortBinding)
			if err != nil {
				continue
			}
			portBindings = append(portBindings, portBinding)
		}
	}

	var autoStartPtr *bool
	if input["autoStart"] != nil {
		var err error
		autoStart, err := sharedHelper.ParseBoolParam(
			input["autoStart"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidAutoStart")
		}
		autoStartPtr = &autoStart
	}

	var timeoutStartSecsPtr *uint
	if input["timeoutStartSecs"] != nil {
		var assertOk bool
		timeoutStartSecs, assertOk := input["timeoutStartSecs"].(uint)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidTimeoutStartSecs")
		}
		timeoutStartSecsPtr = &timeoutStartSecs
	}

	var autoRestartPtr *bool
	if input["autoRestart"] != nil {
		var err error
		autoRestart, err := sharedHelper.ParseBoolParam(
			input["autoRestart"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidAutoRestart")
		}
		autoRestartPtr = &autoRestart
	}

	var maxStartRetriesPtr *uint
	if input["maxStartRetries"] != nil {
		var assertOk bool
		maxStartRetries, assertOk := input["maxStartRetries"].(uint)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidMaxStartRetries")
		}
		maxStartRetriesPtr = &maxStartRetries
	}

	autoCreateMapping := true
	if input["autoCreateMapping"] != nil {
		var err error
		autoCreateMapping, err = sharedHelper.ParseBoolParam(
			input["autoCreateMapping"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidAutoCreateMapping")
		}
	}

	if shouldSchedule {
		cliCmd := infraEnvs.SpeediaOsBinary + " services create-installable"
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
		timeoutSeconds := uint(600)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSeconds, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "CreateInstallableServiceScheduled")
	}

	dto := dto.NewCreateInstallableService(
		name, envs, portBindings, versionPtr, startupFilePtr, autoStartPtr,
		timeoutStartSecsPtr, autoRestartPtr, maxStartRetriesPtr, &autoCreateMapping,
	)

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(service.persistentDbService)
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(service.persistentDbService)
	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(service.persistentDbService)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(service.persistentDbService)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbService)

	err = useCase.CreateInstallableService(
		servicesQueryRepo, servicesCmdRepo, mappingQueryRepo,
		mappingCmdRepo, vhostQueryRepo, dto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "InstallableServiceCreated")
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
		log.Printf("Version: %v", input["version"])
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
			return NewServiceOutput(UserError, "InvalidPortBindings")
		}

		for _, rawPortBinding := range rawPortBindings {
			portBinding, err := valueObject.NewPortBinding(rawPortBinding)
			if err != nil {
				continue
			}
			portBindings = append(portBindings, portBinding)
		}
	}

	autoCreateMapping := true
	if input["autoCreateMapping"] != nil {
		var err error
		autoCreateMapping, err = sharedHelper.ParseBoolParam(
			input["autoCreateMapping"],
		)
		if err != nil {
			return NewServiceOutput(UserError, "InvalidAutoCreateMapping")
		}
	}

	dto := dto.NewCreateCustomService(
		name, svcType, startCmd, []valueObject.ServiceEnv{}, portBindings,
		nil, nil, nil, nil, nil, versionPtr, nil, nil, nil, nil, nil, nil, nil, nil,
		&autoCreateMapping,
	)

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(service.persistentDbService)
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(service.persistentDbService)
	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(service.persistentDbService)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(service.persistentDbService)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbService)

	err = useCase.CreateCustomService(
		servicesQueryRepo, servicesCmdRepo, mappingQueryRepo,
		mappingCmdRepo, vhostQueryRepo, dto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "CustomServiceCreated")
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
		status, err := valueObject.NewServiceStatus(input["status"].(string))
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
			return NewServiceOutput(UserError, "InvalidPortBindings")
		}

		for _, rawPortBinding := range rawPortBindings {
			portBinding, err := valueObject.NewPortBinding(rawPortBinding)
			if err != nil {
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

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(service.persistentDbService)
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(service.persistentDbService)
	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(service.persistentDbService)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(service.persistentDbService)

	err = useCase.UpdateService(
		servicesQueryRepo, servicesCmdRepo, mappingQueryRepo, mappingCmdRepo, dto,
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

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(service.persistentDbService)
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(service.persistentDbService)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(service.persistentDbService)

	err = useCase.DeleteService(
		servicesQueryRepo,
		servicesCmdRepo,
		mappingCmdRepo,
		name,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ServiceDeleted")
}
