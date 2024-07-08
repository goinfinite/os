package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateService struct {
	Name         valueObject.ServiceName     `json:"name"`
	Type         *valueObject.ServiceType    `json:"type"`
	StartCmd     *valueObject.UnixCommand    `json:"startCmd"`
	Status       *valueObject.ServiceStatus  `json:"status"`
	Version      *valueObject.ServiceVersion `json:"version"`
	StartupFile  *valueObject.UnixFilePath   `json:"startupFile"`
	PortBindings []valueObject.PortBinding   `json:"portBindings"`
}

func NewUpdateService(
	name valueObject.ServiceName,
	svcType *valueObject.ServiceType,
	startCmd *valueObject.UnixCommand,
	status *valueObject.ServiceStatus,
	version *valueObject.ServiceVersion,
	startupFile *valueObject.UnixFilePath,
	portBindings []valueObject.PortBinding,
) UpdateService {
	return UpdateService{
		Name:         name,
		Type:         svcType,
		StartCmd:     startCmd,
		Status:       status,
		Version:      version,
		StartupFile:  startupFile,
		PortBindings: portBindings,
	}
}
