package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type InstalledService struct {
	Name                 valueObject.ServiceName              `json:"name"`
	Nature               valueObject.ServiceNature            `json:"nature"`
	Type                 valueObject.ServiceType              `json:"type"`
	Version              valueObject.ServiceVersion           `json:"version"`
	Status               valueObject.ServiceStatus            `json:"status"`
	StartCmd             tkValueObject.UnixCommand            `json:"startCmd"`
	Envs                 []valueObject.ServiceEnv             `json:"envs"`
	PortBindings         []valueObject.PortBinding            `json:"portBindings"`
	StopTimeoutSecs      tkValueObject.UnixTime               `json:"-"`
	StopCmdSteps         []tkValueObject.UnixCommand          `json:"-"`
	PreStartTimeoutSecs  tkValueObject.UnixTime               `json:"-"`
	PreStartCmdSteps     []tkValueObject.UnixCommand          `json:"-"`
	PostStartTimeoutSecs tkValueObject.UnixTime               `json:"-"`
	PostStartCmdSteps    []tkValueObject.UnixCommand          `json:"-"`
	PreStopTimeoutSecs   tkValueObject.UnixTime               `json:"-"`
	PreStopCmdSteps      []tkValueObject.UnixCommand          `json:"-"`
	PostStopTimeoutSecs  tkValueObject.UnixTime               `json:"-"`
	PostStopCmdSteps     []tkValueObject.UnixCommand          `json:"-"`
	ExecUser             *tkValueObject.UnixUsername           `json:"execUser"`
	WorkingDirectory     *tkValueObject.UnixAbsoluteFilePath  `json:"workingDirectory"`
	StartupFile          *tkValueObject.UnixAbsoluteFilePath  `json:"startupFile"`
	AutoStart            *bool                                `json:"autoStart"`
	AutoRestart          *bool                                `json:"autoRestart"`
	TimeoutStartSecs     *uint                                `json:"timeoutStartSecs"`
	MaxStartRetries      *uint                                `json:"maxStartRetries"`
	LogOutputPath        *tkValueObject.UnixAbsoluteFilePath  `json:"logOutputPath"`
	LogErrorPath         *tkValueObject.UnixAbsoluteFilePath  `json:"logErrorPath"`
	AvatarUrl            *tkValueObject.Url                   `json:"avatarUrl"`
	CreatedAt            tkValueObject.UnixTime               `json:"createdAt"`
	UpdatedAt            tkValueObject.UnixTime               `json:"updatedAt"`
}

func NewInstalledService(
	name valueObject.ServiceName,
	nature valueObject.ServiceNature,
	serviceType valueObject.ServiceType,
	version valueObject.ServiceVersion,
	startCmd tkValueObject.UnixCommand,
	status valueObject.ServiceStatus,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	stopTimeoutSecs tkValueObject.UnixTime,
	stopSteps []tkValueObject.UnixCommand,
	preStartTimeoutSecs tkValueObject.UnixTime,
	preStartSteps []tkValueObject.UnixCommand,
	postStartTimeoutSecs tkValueObject.UnixTime,
	postStartSteps []tkValueObject.UnixCommand,
	preStopTimeoutSecs tkValueObject.UnixTime,
	preStopSteps []tkValueObject.UnixCommand,
	postStopTimeoutSecs tkValueObject.UnixTime,
	postStopSteps []tkValueObject.UnixCommand,
	execUser *tkValueObject.UnixUsername,
	workingDirectory, startupFile *tkValueObject.UnixAbsoluteFilePath,
	autoStart, autoRestart *bool,
	timeoutStartSecs, maxStartRetries *uint,
	logOutputPath, logErrorPath *tkValueObject.UnixAbsoluteFilePath,
	avatarUrl *tkValueObject.Url,
	createdAt tkValueObject.UnixTime,
	updatedAt tkValueObject.UnixTime,
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
