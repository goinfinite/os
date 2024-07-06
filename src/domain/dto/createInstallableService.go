package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateInstallableService struct {
	Name              valueObject.ServiceName     `json:"name"`
	Envs              []valueObject.ServiceEnv    `json:"envs"`
	PortBindings      []valueObject.PortBinding   `json:"portBindings"`
	AutoCreateMapping bool                        `json:"autoCreateMapping"`
	Version           *valueObject.ServiceVersion `json:"version"`
	StartupFile       *valueObject.UnixFilePath   `json:"startupFile"`
}

func NewCreateInstallableService(
	name valueObject.ServiceName,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	autoCreateMapping bool,
	version *valueObject.ServiceVersion,
	startupFile *valueObject.UnixFilePath,
) CreateInstallableService {
	return CreateInstallableService{
		Name:              name,
		Envs:              envs,
		PortBindings:      portBindings,
		AutoCreateMapping: autoCreateMapping,
		Version:           version,
		StartupFile:       startupFile,
	}
}
