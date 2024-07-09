package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateService struct {
	Name              valueObject.ServiceName     `json:"name"`
	Type              *valueObject.ServiceType    `json:"type"`
	StartCmd          *valueObject.UnixCommand    `json:"startCmd"`
	Status            *valueObject.ServiceStatus  `json:"status"`
	Version           *valueObject.ServiceVersion `json:"version"`
	StartupFile       *valueObject.UnixFilePath   `json:"startupFile"`
	Envs              []valueObject.ServiceEnv    `json:"envs"`
	PortBindings      []valueObject.PortBinding   `json:"portBindings"`
	AutoStart         *bool                       `json:"autoStart"`
	TimeoutStartSecs  *uint                       `json:"timeoutStartSecs"`
	AutoRestart       *bool                       `json:"autoRestart"`
	MaxStartRetries   *uint                       `json:"maxStartRetries"`
	AutoCreateMapping *bool                       `json:"autoCreateMapping"`
}

func NewUpdateService(
	name valueObject.ServiceName,
	svcType *valueObject.ServiceType,
	startCmd *valueObject.UnixCommand,
	status *valueObject.ServiceStatus,
	version *valueObject.ServiceVersion,
	startupFile *valueObject.UnixFilePath,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	autoStart *bool,
	timeoutStartSecs *uint,
	autoRestart *bool,
	maxStartRetries *uint,
) UpdateService {
	return UpdateService{
		Name:             name,
		Type:             svcType,
		StartCmd:         startCmd,
		Status:           status,
		Version:          version,
		StartupFile:      startupFile,
		Envs:             envs,
		PortBindings:     portBindings,
		AutoStart:        autoStart,
		TimeoutStartSecs: timeoutStartSecs,
		AutoRestart:      autoRestart,
		MaxStartRetries:  maxStartRetries,
	}
}
