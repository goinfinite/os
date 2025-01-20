package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type InstalledService struct {
	Name              valueObject.ServiceName    `json:"name"`
	Nature            valueObject.ServiceNature  `json:"nature"`
	Type              valueObject.ServiceType    `json:"type"`
	Version           valueObject.ServiceVersion `json:"version"`
	Status            valueObject.ServiceStatus  `json:"status"`
	StartCmd          valueObject.UnixCommand    `json:"startCmd"`
	AvatarUrl         *valueObject.Url           `json:"avatarUrl"`
	Envs              []valueObject.ServiceEnv   `json:"envs"`
	PortBindings      []valueObject.PortBinding  `json:"portBindings"`
	StopCmdSteps      []valueObject.UnixCommand  `json:"stopCmdSteps"`
	PreStartCmdSteps  []valueObject.UnixCommand  `json:"preStartCmdSteps"`
	PostStartCmdSteps []valueObject.UnixCommand  `json:"postStartCmdSteps"`
	PreStopCmdSteps   []valueObject.UnixCommand  `json:"preStopCmdSteps"`
	PostStopCmdSteps  []valueObject.UnixCommand  `json:"postStopCmdSteps"`
	ExecUser          *valueObject.UnixUsername  `json:"execUser"`
	WorkingDirectory  *valueObject.UnixFilePath  `json:"workingDirectory"`
	StartupFile       *valueObject.UnixFilePath  `json:"startupFile"`
	AutoStart         *bool                      `json:"autoStart"`
	AutoRestart       *bool                      `json:"autoRestart"`
	TimeoutStartSecs  *uint                      `json:"timeoutStartSecs"`
	MaxStartRetries   *uint                      `json:"maxStartRetries"`
	LogOutputPath     *valueObject.UnixFilePath  `json:"logOutputPath"`
	LogErrorPath      *valueObject.UnixFilePath  `json:"logErrorPath"`
	CreatedAt         valueObject.UnixTime       `json:"createdAt"`
	UpdatedAt         valueObject.UnixTime       `json:"updatedAt"`
}

func NewInstalledService(
	name valueObject.ServiceName,
	nature valueObject.ServiceNature,
	serviceType valueObject.ServiceType,
	version valueObject.ServiceVersion,
	startCmd valueObject.UnixCommand,
	status valueObject.ServiceStatus,
	avatarUrl *valueObject.Url,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	stopSteps, preStartSteps, postStartSteps, preStopSteps, postStopSteps []valueObject.UnixCommand,
	execUser *valueObject.UnixUsername,
	workingDirectory, startupFile *valueObject.UnixFilePath,
	autoStart, autoRestart *bool,
	timeoutStartSecs, maxStartRetries *uint,
	logOutputPath, logErrorPath *valueObject.UnixFilePath,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) InstalledService {
	return InstalledService{
		Name:              name,
		Nature:            nature,
		Type:              serviceType,
		Version:           version,
		StartCmd:          startCmd,
		Status:            status,
		AvatarUrl:         avatarUrl,
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
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}
