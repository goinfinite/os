package liaison

import (
	"errors"
	"strconv"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	domainRepository "github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	runtimeInfra "github.com/goinfinite/os/src/infra/runtime"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	liaisonHelper "github.com/goinfinite/os/src/presentation/liaison/helper"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
	tkInfra "github.com/goinfinite/tk/src/infra"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
)

const phpModuleUpdateTaskTimeoutSecs uint16 = 1800

type RuntimeLiaison struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	availabilityInspector *liaisonHelper.ServiceAvailabilityInspector
	runtimeQueryRepo      *runtimeInfra.RuntimeQueryRepo
	runtimeCmdRepo        *runtimeInfra.RuntimeCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
	phpServiceName        valueObject.ServiceName
}

func NewRuntimeLiaison(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *RuntimeLiaison {
	return &RuntimeLiaison{
		persistentDbSvc: persistentDbSvc,
		availabilityInspector: liaisonHelper.NewServiceAvailabilityInspector(
			persistentDbSvc,
		),
		runtimeQueryRepo:      runtimeInfra.NewRuntimeQueryRepo(),
		runtimeCmdRepo:        runtimeInfra.NewRuntimeCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
		phpServiceName:        valueObject.ServiceNamePhpWebServer,
	}
}

func (liaison *RuntimeLiaison) ReadPhpConfigs(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	if !liaison.availabilityInspector.IsAvailable(liaison.phpServiceName) {
		return liaisonHelper.NewServiceUnavailableResponse()
	}

	hostname, err := tkValueObject.NewFqdn(untrustedInput["hostname"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	phpConfigs, err := useCase.ReadPhpConfigs(liaison.runtimeQueryRepo, hostname)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, phpConfigs,
	)
}

func (liaison *RuntimeLiaison) phpModulesUpdateRequestFactory(
	hostname tkValueObject.Fqdn,
	phpVersion valueObject.PhpVersion,
	modules []entity.PhpModule,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) dto.UpdatePhpModulesRequest {
	moduleUpdates := []dto.PhpModuleUpdate{}
	for _, module := range modules {
		moduleUpdates = append(
			moduleUpdates,
			dto.NewPhpModuleUpdate(module.Name, module.Status),
		)
	}

	return dto.NewUpdatePhpModulesRequest(
		hostname, phpVersion, moduleUpdates,
		operatorAccountId, operatorIpAddress,
	)
}

func (liaison *RuntimeLiaison) executePhpModulesUpdate(
	requestDto dto.UpdatePhpModulesRequest,
	vhostQueryRepo *vhostInfra.VirtualHostQueryRepo,
) tkPresentation.LiaisonResponse {
	updatePhpModulesUseCase := useCase.NewUpdatePhpModules(
		liaison.runtimeQueryRepo,
		liaison.runtimeCmdRepo,
		vhostQueryRepo,
		liaison.activityRecordCmdRepo,
	)
	responseDto, err := updatePhpModulesUseCase.Execute(requestDto)
	if err != nil {
		responseStatus := tkPresentation.LiaisonResponseStatusInfraError
		if errors.Is(err, domainRepository.ErrPhpModuleNotSupportedForVersion) {
			responseStatus = tkPresentation.LiaisonResponseStatusUserError
		}
		if len(responseDto.ModulesSuccessfullyUpdated) > 0 ||
			len(responseDto.FailedModulesWithReason) > 0 {
			return tkPresentation.NewLiaisonResponse(
				responseStatus, responseDto, err.Error(),
			)
		}
		return tkPresentation.NewLiaisonResponseNoMessage(
			responseStatus, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, responseDto,
	)
}

func (liaison *RuntimeLiaison) UpdatePhpModules(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	if !liaison.availabilityInspector.IsAvailable(liaison.phpServiceName) {
		return liaisonHelper.NewServiceUnavailableResponse()
	}

	requiredParams := []string{"hostname", "version", "modules"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	hostname, err := tkValueObject.NewFqdn(untrustedInput["hostname"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}
	phpVersion, err := valueObject.NewPhpVersion(untrustedInput["version"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	phpModules, assertOk := untrustedInput["modules"].([]entity.PhpModule)
	if !assertOk || len(phpModules) == 0 {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, "InvalidPhpModules",
		)
	}

	operatorAccountId, operatorIpAddress, err := liaisonHelper.ReadOperatorContext(untrustedInput)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	moduleUpdateRequestDto := liaison.phpModulesUpdateRequestFactory(
		hostname, phpVersion, phpModules, operatorAccountId, operatorIpAddress,
	)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(liaison.persistentDbSvc)
	return liaison.executePhpModulesUpdate(moduleUpdateRequestDto, vhostQueryRepo)
}

func (liaison *RuntimeLiaison) RunPhpCommand(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	if !liaison.availabilityInspector.IsAvailable(liaison.phpServiceName) {
		return liaisonHelper.NewServiceUnavailableResponse()
	}

	requiredParams := []string{"hostname", "command"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	hostname, err := tkValueObject.NewFqdn(untrustedInput["hostname"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	command, err := tkValueObject.NewUnixCommand(untrustedInput["command"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	var timeoutSecsPtr *uint64
	if untrustedInput["timeoutSecs"] != nil {
		timeoutSecs, err := tkVoUtil.InterfaceToUint64(untrustedInput["timeoutSecs"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, "TimeoutSecsMustBeUint64",
			)
		}
		timeoutSecsPtr = &timeoutSecs
	}

	operatorAccountId, operatorIpAddress, err := liaisonHelper.ReadOperatorContext(untrustedInput)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	runRequestDto := dto.NewRunPhpCommandRequest(
		hostname, command, timeoutSecsPtr, operatorAccountId, operatorIpAddress,
	)

	accountQueryRepo := accountInfra.NewAccountQueryRepo(liaison.persistentDbSvc)

	runResponse, err := useCase.RunPhpCommand(
		accountQueryRepo, liaison.runtimeCmdRepo, runRequestDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, runResponse,
	)
}

func (liaison *RuntimeLiaison) readPhpConfigsRequest(
	untrustedInput map[string]any,
) (dto.UpdatePhpConfigsRequest, []dto.PhpModuleParsingFailure, error) {
	requiredParams := []string{"hostname", "version"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return dto.UpdatePhpConfigsRequest{}, nil, err
	}

	hostname, err := tkValueObject.NewFqdn(untrustedInput["hostname"])
	if err != nil {
		return dto.UpdatePhpConfigsRequest{}, nil, err
	}

	phpVersion, err := valueObject.NewPhpVersion(untrustedInput["version"])
	if err != nil {
		return dto.UpdatePhpConfigsRequest{}, nil, err
	}

	phpModules := []entity.PhpModule{}
	if _, exists := untrustedInput["modules"]; exists {
		var assertOk bool
		phpModules, assertOk = untrustedInput["modules"].([]entity.PhpModule)
		if !assertOk {
			return dto.UpdatePhpConfigsRequest{}, nil,
				errors.New("InvalidPhpModules")
		}
	}

	moduleParsingFailures := []dto.PhpModuleParsingFailure{}
	if _, exists := untrustedInput["moduleParsingFailures"]; exists {
		var assertOk bool
		rawFailures := untrustedInput["moduleParsingFailures"]
		moduleParsingFailures, assertOk = rawFailures.([]dto.PhpModuleParsingFailure)
		if !assertOk {
			return dto.UpdatePhpConfigsRequest{}, nil,
				errors.New("InvalidPhpModuleParsingFailures")
		}
	}

	phpSettings := []entity.PhpSetting{}
	if _, exists := untrustedInput["settings"]; exists {
		var assertOk bool
		phpSettings, assertOk = untrustedInput["settings"].([]entity.PhpSetting)
		if !assertOk {
			return dto.UpdatePhpConfigsRequest{}, nil,
				errors.New("InvalidPhpSettings")
		}
	}

	operatorAccountId, operatorIpAddress, err := liaisonHelper.ReadOperatorContext(untrustedInput)
	if err != nil {
		return dto.UpdatePhpConfigsRequest{}, nil, err
	}

	requestDto := dto.NewUpdatePhpConfigsRequest(
		hostname, phpVersion, phpModules, phpSettings,
		operatorAccountId, operatorIpAddress,
	)
	return requestDto, moduleParsingFailures, nil
}

func (liaison *RuntimeLiaison) updatePhpVersion(
	requestDto dto.UpdatePhpConfigsRequest,
	vhostQueryRepo *vhostInfra.VirtualHostQueryRepo,
) error {
	updateVersionRequestDto := dto.NewUpdatePhpVersionRequest(
		requestDto.Hostname, requestDto.PhpVersion,
		requestDto.OperatorAccountId, requestDto.OperatorIpAddress,
	)
	updateVersionUseCase := useCase.NewUpdatePhpVersion(
		liaison.runtimeQueryRepo,
		liaison.runtimeCmdRepo,
		vhostQueryRepo,
		liaison.activityRecordCmdRepo,
	)
	return updateVersionUseCase.Execute(updateVersionRequestDto)
}

func (liaison *RuntimeLiaison) updatePhpSettings(
	requestDto dto.UpdatePhpConfigsRequest,
	vhostQueryRepo *vhostInfra.VirtualHostQueryRepo,
) error {
	updateSettingsRequestDto := dto.NewUpdatePhpSettingsRequest(
		requestDto.Hostname, requestDto.PhpSettings,
		requestDto.OperatorAccountId, requestDto.OperatorIpAddress,
	)
	updateSettingsUseCase := useCase.NewUpdatePhpSettings(
		liaison.runtimeCmdRepo,
		vhostQueryRepo,
		liaison.activityRecordCmdRepo,
	)
	return updateSettingsUseCase.Execute(updateSettingsRequestDto)
}

func (liaison *RuntimeLiaison) phpModuleUpdateCliCommandFactory(
	requestDto dto.UpdatePhpModulesRequest,
) (tkValueObject.UnixCommand, error) {
	commandParts := []string{
		infraEnvs.InfiniteOsBinary,
		"runtime",
		"php",
		"update-modules",
		"--hostname",
		tkInfra.ShellEscape{}.Quote(requestDto.Hostname.String()),
		"--version",
		tkInfra.ShellEscape{}.Quote(requestDto.PhpVersion.String()),
		"--operator-account-id",
		tkInfra.ShellEscape{}.Quote(requestDto.OperatorAccountId.String()),
		"--operator-ip-address",
		tkInfra.ShellEscape{}.Quote(requestDto.OperatorIpAddress.String()),
	}
	for _, module := range requestDto.PhpModules {
		moduleParam := module.Name.String() + ":" + strconv.FormatBool(module.Status)
		commandParts = append(
			commandParts,
			"--module",
			tkInfra.ShellEscape{}.Quote(moduleParam),
		)
	}

	return tkValueObject.NewUnixCommand(strings.Join(commandParts, " "))
}

func (liaison *RuntimeLiaison) schedulePhpModulesUpdate(
	requestDto dto.UpdatePhpModulesRequest,
) (valueObject.ScheduledTaskId, error) {
	command, err := liaison.phpModuleUpdateCliCommandFactory(requestDto)
	if err != nil {
		return 0, err
	}

	taskName, err := valueObject.NewScheduledTaskName("UpdatePhpModules")
	if err != nil {
		return 0, err
	}
	taskTag, err := valueObject.NewScheduledTaskTag("runtime")
	if err != nil {
		return 0, err
	}
	timeoutSecs := phpModuleUpdateTaskTimeoutSecs
	createDto := dto.NewCreateScheduledTask(
		taskName, command, []valueObject.ScheduledTaskTag{taskTag}, &timeoutSecs, nil,
	)
	scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(
		liaison.persistentDbSvc,
	)

	return useCase.CreateScheduledTask(scheduledTaskCmdRepo, createDto)
}

func (liaison *RuntimeLiaison) UpdatePhpConfigs(
	untrustedInput map[string]any,
	shouldScheduleModuleUpdates bool,
) tkPresentation.LiaisonResponse {
	if !liaison.availabilityInspector.IsAvailable(liaison.phpServiceName) {
		return liaisonHelper.NewServiceUnavailableResponse()
	}

	requestDto, moduleParsingFailures, readErr := liaison.readPhpConfigsRequest(
		untrustedInput,
	)
	if readErr != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, readErr.Error(),
		)
	}

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(liaison.persistentDbSvc)

	err := liaison.updatePhpVersion(requestDto, vhostQueryRepo)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	if len(requestDto.PhpSettings) > 0 {
		err = liaison.updatePhpSettings(requestDto, vhostQueryRepo)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
			)
		}
	}

	if len(requestDto.PhpModules) == 0 {
		if len(moduleParsingFailures) > 0 {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusMultiStatus,
				dto.NewUpdatePhpConfigsResponse(nil, moduleParsingFailures),
			)
		}

		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusSuccess, "PhpConfigsUpdated",
		)
	}

	moduleUpdateRequestDto := liaison.phpModulesUpdateRequestFactory(
		requestDto.Hostname, requestDto.PhpVersion, requestDto.PhpModules,
		requestDto.OperatorAccountId, requestDto.OperatorIpAddress,
	)
	if shouldScheduleModuleUpdates {
		taskId, scheduleErr := liaison.schedulePhpModulesUpdate(moduleUpdateRequestDto)
		if scheduleErr != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusInfraError,
				scheduleErr.Error(),
			)
		}

		responseStatus := tkPresentation.LiaisonResponseStatusCreated
		if len(moduleParsingFailures) > 0 {
			responseStatus = tkPresentation.LiaisonResponseStatusMultiStatus
		}
		taskResponseDto := dto.NewUpdatePhpConfigsResponse(
			&taskId, moduleParsingFailures,
		)

		return tkPresentation.NewLiaisonResponseNoMessage(
			responseStatus, taskResponseDto,
		)
	}

	return liaison.executePhpModulesUpdate(moduleUpdateRequestDto, vhostQueryRepo)
}
