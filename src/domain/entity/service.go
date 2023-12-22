package entity

import "github.com/speedianet/os/src/domain/valueObject"

type Service struct {
	Name         valueObject.ServiceName    `json:"name"`
	Nature       valueObject.ServiceNature  `json:"nature"`
	Type         valueObject.ServiceType    `json:"type"`
	Version      valueObject.ServiceVersion `json:"version"`
	Command      valueObject.UnixCommand    `json:"command"`
	Status       valueObject.ServiceStatus  `json:"status"`
	StartupFile  *valueObject.UnixFilePath  `json:"startupFile,omitempty"`
	PortBindings []valueObject.PortBinding  `json:"portBindings,omitempty"`
}

func NewService(
	name valueObject.ServiceName,
	nature valueObject.ServiceNature,
	svcType valueObject.ServiceType,
	version valueObject.ServiceVersion,
	command valueObject.UnixCommand,
	status valueObject.ServiceStatus,
	startupFile *valueObject.UnixFilePath,
	portBindings []valueObject.PortBinding,
) Service {
	return Service{
		Name:         name,
		Nature:       nature,
		Type:         svcType,
		Version:      version,
		Command:      command,
		Status:       status,
		StartupFile:  startupFile,
		PortBindings: portBindings,
	}
}
