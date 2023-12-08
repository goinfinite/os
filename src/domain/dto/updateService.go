package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateService struct {
	Name        valueObject.ServiceName     `json:"name"`
	Type        *valueObject.ServiceType    `json:"type,omitempty"`
	Command     *valueObject.UnixCommand    `json:"command,omitempty"`
	Status      *valueObject.ServiceStatus  `json:"status,omitempty"`
	Version     *valueObject.ServiceVersion `json:"version,omitempty"`
	StartupFile *valueObject.UnixFilePath   `json:"startupFile,omitempty"`
	Ports       []valueObject.NetworkPort   `json:"ports,omitempty"`
}

func NewUpdateService(
	name valueObject.ServiceName,
	svcType *valueObject.ServiceType,
	command *valueObject.UnixCommand,
	status *valueObject.ServiceStatus,
	version *valueObject.ServiceVersion,
	startupFile *valueObject.UnixFilePath,
	ports []valueObject.NetworkPort,
) UpdateService {
	return UpdateService{
		Name:        name,
		Type:        svcType,
		Command:     command,
		Status:      status,
		Version:     version,
		StartupFile: startupFile,
		Ports:       ports,
	}
}
