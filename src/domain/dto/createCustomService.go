package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateCustomService struct {
	Name              valueObject.ServiceName     `json:"name"`
	Type              valueObject.ServiceType     `json:"type"`
	StartCmd          valueObject.UnixCommand     `json:"startCmd"`
	Envs              []valueObject.ServiceEnv    `json:"envs"`
	PortBindings      []valueObject.PortBinding   `json:"portBindings"`
	Version           *valueObject.ServiceVersion `json:"version"`
	AutoStart         *bool                       `json:"autoStart"`
	TimeoutStartSecs  *uint                       `json:"timeoutStartSecs"`
	AutoRestart       *bool                       `json:"autoRestart"`
	MaxStartRetries   *uint                       `json:"maxStartRetries"`
	AutoCreateMapping *bool                       `json:"autoCreateMapping"`
}

func NewCreateCustomService(
	name valueObject.ServiceName,
	serviceType valueObject.ServiceType,
	startCmd valueObject.UnixCommand,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	version *valueObject.ServiceVersion,
	autoStart *bool,
	timeoutStartSecs *uint,
	autoRestart *bool,
	maxStartRetries *uint,
	autoCreateMapping *bool,
) CreateCustomService {
	return CreateCustomService{
		Name:              name,
		Type:              serviceType,
		StartCmd:          startCmd,
		Envs:              envs,
		PortBindings:      portBindings,
		Version:           version,
		AutoStart:         autoStart,
		TimeoutStartSecs:  timeoutStartSecs,
		AutoRestart:       autoRestart,
		MaxStartRetries:   maxStartRetries,
		AutoCreateMapping: autoCreateMapping,
	}
}
