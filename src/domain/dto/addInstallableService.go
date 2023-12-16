package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddInstallableService struct {
	Name         valueObject.ServiceName     `json:"name"`
	Version      *valueObject.ServiceVersion `json:"version,omitempty"`
	StartupFile  *valueObject.UnixFilePath   `json:"startupFile,omitempty"`
	PortBindings []valueObject.PortBinding   `json:"portBindings,omitempty"`
}

func NewAddInstallableService(
	name valueObject.ServiceName,
	version *valueObject.ServiceVersion,
	startupFile *valueObject.UnixFilePath,
	portBindings []valueObject.PortBinding,
) AddInstallableService {
	return AddInstallableService{
		Name:         name,
		Version:      version,
		StartupFile:  startupFile,
		PortBindings: portBindings,
	}
}
