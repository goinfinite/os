package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateService struct {
	Name              valueObject.ServiceName     `json:"name"`
	Type              *valueObject.ServiceType    `json:"type"`
	Version           *valueObject.ServiceVersion `json:"version"`
	Status            *valueObject.ServiceStatus  `json:"status"`
	StartCmd          *valueObject.UnixCommand    `json:"startCmd"`
	Envs              []valueObject.ServiceEnv    `json:"envs"`
	PortBindings      []valueObject.PortBinding   `json:"portBindings"`
	StopCmdSteps      []valueObject.UnixCommand   `json:"stopCmdSteps"`
	PreStartCmdSteps  []valueObject.UnixCommand   `json:"preStartCmdSteps"`
	PostStartCmdSteps []valueObject.UnixCommand   `json:"postStartCmdSteps"`
	PreStopCmdSteps   []valueObject.UnixCommand   `json:"preStopCmdSteps"`
	PostStopCmdSteps  []valueObject.UnixCommand   `json:"postStopCmdSteps"`
	ExecUser          *valueObject.UnixUsername   `json:"execUser"`
	WorkingDirectory  *valueObject.UnixFilePath   `json:"workingDirectory"`
	StartupFile       *valueObject.UnixFilePath   `json:"startupFile"`
	AutoStart         *bool                       `json:"autoStart"`
	AutoRestart       *bool                       `json:"autoRestart"`
	TimeoutStartSecs  *uint                       `json:"timeoutStartSecs"`
	MaxStartRetries   *uint                       `json:"maxStartRetries"`
	LogOutputPath     *valueObject.UnixFilePath   `json:"logOutputPath"`
	LogErrorPath      *valueObject.UnixFilePath   `json:"logErrorPath"`
}

func NewUpdateService(
	name valueObject.ServiceName,
	svcType *valueObject.ServiceType,
	version *valueObject.ServiceVersion,
	status *valueObject.ServiceStatus,
	startCmd *valueObject.UnixCommand,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	stopSteps, preStartSteps, postStartSteps, preStopSteps, postStopSteps []valueObject.UnixCommand,
	execUser *valueObject.UnixUsername,
	workingDirectory, startupFile *valueObject.UnixFilePath,
	autoStart, autoRestart *bool,
	timeoutStartSecs, maxStartRetries *uint,
	logOutputPath, logErrorPath *valueObject.UnixFilePath,
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
	}
}
