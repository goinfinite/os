package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type InstalledService struct {
	Name                 valueObject.ServiceName    `json:"name"`
	Nature               valueObject.ServiceNature  `json:"nature"`
	Type                 valueObject.ServiceType    `json:"type"`
	Version              valueObject.ServiceVersion `json:"version"`
	Status               valueObject.ServiceStatus  `json:"status"`
	StartCmd             valueObject.UnixCommand    `json:"startCmd"`
	Envs                 []valueObject.ServiceEnv   `json:"envs"`
	PortBindings         []valueObject.PortBinding  `json:"portBindings"`
	StopTimeoutSecs      valueObject.UnixTime       `json:"-"`
	StopCmdSteps         []valueObject.UnixCommand  `json:"-"`
	PreStartTimeoutSecs  valueObject.UnixTime       `json:"-"`
	PreStartCmdSteps     []valueObject.UnixCommand  `json:"-"`
	PostStartTimeoutSecs valueObject.UnixTime       `json:"-"`
	PostStartCmdSteps    []valueObject.UnixCommand  `json:"-"`
	PreStopTimeoutSecs   valueObject.UnixTime       `json:"-"`
	PreStopCmdSteps      []valueObject.UnixCommand  `json:"-"`
	PostStopTimeoutSecs  valueObject.UnixTime       `json:"-"`
	PostStopCmdSteps     []valueObject.UnixCommand  `json:"-"`
	ExecUser             *valueObject.UnixUsername  `json:"execUser"`
	WorkingDirectory     *valueObject.UnixFilePath  `json:"workingDirectory"`
	StartupFile          *valueObject.UnixFilePath  `json:"startupFile"`
	AutoStart            *bool                      `json:"autoStart"`
	AutoRestart          *bool                      `json:"autoRestart"`
	TimeoutStartSecs     *uint                      `json:"timeoutStartSecs"`
	MaxStartRetries      *uint                      `json:"maxStartRetries"`
	LogOutputPath        *valueObject.UnixFilePath  `json:"logOutputPath"`
	LogErrorPath         *valueObject.UnixFilePath  `json:"logErrorPath"`
	AvatarUrl            *valueObject.Url           `json:"avatarUrl"`
	CreatedAt            valueObject.UnixTime       `json:"createdAt"`
	UpdatedAt            valueObject.UnixTime       `json:"updatedAt"`
}

func NewInstalledService(
	name valueObject.ServiceName,
	nature valueObject.ServiceNature,
	serviceType valueObject.ServiceType,
	version valueObject.ServiceVersion,
	startCmd valueObject.UnixCommand,
	status valueObject.ServiceStatus,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	stopTimeoutSecs valueObject.UnixTime,
	stopSteps []valueObject.UnixCommand,
	preStartTimeoutSecs valueObject.UnixTime,
	preStartSteps []valueObject.UnixCommand,
	postStartTimeoutSecs valueObject.UnixTime,
	postStartSteps []valueObject.UnixCommand,
	preStopTimeoutSecs valueObject.UnixTime,
	preStopSteps []valueObject.UnixCommand,
	postStopTimeoutSecs valueObject.UnixTime,
	postStopSteps []valueObject.UnixCommand,
	execUser *valueObject.UnixUsername,
	workingDirectory, startupFile *valueObject.UnixFilePath,
	autoStart, autoRestart *bool,
	timeoutStartSecs, maxStartRetries *uint,
	logOutputPath, logErrorPath *valueObject.UnixFilePath,
	avatarUrl *valueObject.Url,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) InstalledService {
	return InstalledService{
		Name:                 name,
		Nature:               nature,
		Type:                 serviceType,
		Version:              version,
		StartCmd:             startCmd,
		Status:               status,
		Envs:                 envs,
		PortBindings:         portBindings,
		StopTimeoutSecs:      stopTimeoutSecs,
		StopCmdSteps:         stopSteps,
		PreStartTimeoutSecs:  preStartTimeoutSecs,
		PreStartCmdSteps:     preStartSteps,
		PostStartTimeoutSecs: postStartTimeoutSecs,
		PostStartCmdSteps:    postStartSteps,
		PreStopTimeoutSecs:   preStopTimeoutSecs,
		PreStopCmdSteps:      preStopSteps,
		PostStopTimeoutSecs:  postStopTimeoutSecs,
		PostStopCmdSteps:     postStopSteps,
		ExecUser:             execUser,
		WorkingDirectory:     workingDirectory,
		StartupFile:          startupFile,
		AutoStart:            autoStart,
		AutoRestart:          autoRestart,
		TimeoutStartSecs:     timeoutStartSecs,
		MaxStartRetries:      maxStartRetries,
		LogOutputPath:        logOutputPath,
		LogErrorPath:         logErrorPath,
		AvatarUrl:            avatarUrl,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
}
