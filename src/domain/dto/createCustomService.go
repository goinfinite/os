package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateCustomService struct {
	Name              valueObject.ServiceName     `json:"name"`
	Type              valueObject.ServiceType     `json:"type"`
	Command           valueObject.UnixCommand     `json:"command"`
	Envs              []valueObject.ServiceEnv    `json:"envs"`
	PortBindings      []valueObject.PortBinding   `json:"portBindings"`
	AutoCreateMapping bool                        `json:"autoCreateMapping"`
	Version           *valueObject.ServiceVersion `json:"version"`
}

func NewCreateCustomService(
	name valueObject.ServiceName,
	serviceType valueObject.ServiceType,
	command valueObject.UnixCommand,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	autoCreateMapping bool,
	version *valueObject.ServiceVersion,
) CreateCustomService {
	return CreateCustomService{
		Name:              name,
		Type:              serviceType,
		Command:           command,
		Envs:              envs,
		PortBindings:      portBindings,
		AutoCreateMapping: autoCreateMapping,
		Version:           version,
	}
}
