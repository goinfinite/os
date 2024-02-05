package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddInstallableService struct {
	Name              valueObject.ServiceName     `json:"name"`
	Version           *valueObject.ServiceVersion `json:"version"`
	StartupFile       *valueObject.UnixFilePath   `json:"startupFile"`
	PortBindings      []valueObject.PortBinding   `json:"portBindings"`
	AutoCreateMapping bool                        `json:"autoCreateMapping"`
}

func NewAddInstallableService(
	name valueObject.ServiceName,
	version *valueObject.ServiceVersion,
	startupFile *valueObject.UnixFilePath,
	portBindings []valueObject.PortBinding,
	autoCreateMapping bool,
) AddInstallableService {
	return AddInstallableService{
		Name:              name,
		Version:           version,
		StartupFile:       startupFile,
		PortBindings:      portBindings,
		AutoCreateMapping: autoCreateMapping,
	}
}
