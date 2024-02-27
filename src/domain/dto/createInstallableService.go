package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateInstallableService struct {
	Name              valueObject.ServiceName     `json:"name"`
	Version           *valueObject.ServiceVersion `json:"version"`
	StartupFile       *valueObject.UnixFilePath   `json:"startupFile"`
	PortBindings      []valueObject.PortBinding   `json:"portBindings"`
	AutoCreateMapping bool                        `json:"autoCreateMapping"`
}

func NewCreateInstallableService(
	name valueObject.ServiceName,
	version *valueObject.ServiceVersion,
	startupFile *valueObject.UnixFilePath,
	portBindings []valueObject.PortBinding,
	autoCreateMapping bool,
) CreateInstallableService {
	return CreateInstallableService{
		Name:              name,
		Version:           version,
		StartupFile:       startupFile,
		PortBindings:      portBindings,
		AutoCreateMapping: autoCreateMapping,
	}
}
