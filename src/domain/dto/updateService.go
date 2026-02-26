package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdateService struct {
	Name              valueObject.ServiceName              `json:"name"`
	Type              *valueObject.ServiceType             `json:"type"`
	Version           *valueObject.ServiceVersion          `json:"version"`
	Status            *valueObject.ServiceStatus           `json:"status"`
	StartCmd          *tkValueObject.UnixCommand           `json:"startCmd"`
	Envs              []valueObject.ServiceEnv             `json:"envs"`
	PortBindings      []valueObject.PortBinding            `json:"portBindings"`
	StopCmdSteps      []tkValueObject.UnixCommand          `json:"stopCmdSteps"`
	PreStartCmdSteps  []tkValueObject.UnixCommand          `json:"preStartCmdSteps"`
	PostStartCmdSteps []tkValueObject.UnixCommand          `json:"postStartCmdSteps"`
	PreStopCmdSteps   []tkValueObject.UnixCommand          `json:"preStopCmdSteps"`
	PostStopCmdSteps  []tkValueObject.UnixCommand          `json:"postStopCmdSteps"`
	ExecUser          *tkValueObject.UnixUsername           `json:"execUser"`
	WorkingDirectory  *tkValueObject.UnixAbsoluteFilePath  `json:"workingDirectory"`
	StartupFile       *tkValueObject.UnixAbsoluteFilePath  `json:"startupFile"`
	AutoStart         *bool                                `json:"autoStart"`
	AutoRestart       *bool                                `json:"autoRestart"`
	TimeoutStartSecs  *uint                                `json:"timeoutStartSecs"`
	MaxStartRetries   *uint                                `json:"maxStartRetries"`
	LogOutputPath     *tkValueObject.UnixAbsoluteFilePath  `json:"logOutputPath"`
	LogErrorPath      *tkValueObject.UnixAbsoluteFilePath  `json:"logErrorPath"`
	AvatarUrl         *tkValueObject.Url                   `json:"avatarUrl"`
	OperatorAccountId tkValueObject.AccountId              `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress              `json:"-"`
}

func NewUpdateService(
	name valueObject.ServiceName,
	svcType *valueObject.ServiceType,
	version *valueObject.ServiceVersion,
	status *valueObject.ServiceStatus,
	startCmd *tkValueObject.UnixCommand,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	stopSteps, preStartSteps, postStartSteps, preStopSteps, postStopSteps []tkValueObject.UnixCommand,
	execUser *tkValueObject.UnixUsername,
	workingDirectory, startupFile *tkValueObject.UnixAbsoluteFilePath,
	autoStart, autoRestart *bool,
	timeoutStartSecs, maxStartRetries *uint,
	logOutputPath, logErrorPath *tkValueObject.UnixAbsoluteFilePath,
	avatarUrl *tkValueObject.Url,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) UpdateService {
	return UpdateService{
		Name:              name,
		Type:              svcType,
		Version:           version,
		Status:            status,
		StartCmd:          startCmd,
		Envs:              envs,
		PortBindings:      portBindings,
		StopCmdSteps:      stopSteps,
		PreStartCmdSteps:  preStartSteps,
		PostStartCmdSteps: postStartSteps,
		PreStopCmdSteps:   preStopSteps,
		PostStopCmdSteps:  postStopSteps,
		ExecUser:          execUser,
		WorkingDirectory:  workingDirectory,
		StartupFile:       startupFile,
		AutoStart:         autoStart,
		AutoRestart:       autoRestart,
		TimeoutStartSecs:  timeoutStartSecs,
		MaxStartRetries:   maxStartRetries,
		LogOutputPath:     logOutputPath,
		LogErrorPath:      logErrorPath,
		AvatarUrl:         avatarUrl,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
