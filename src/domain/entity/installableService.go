package entity

import "github.com/speedianet/os/src/domain/valueObject"

type InstallableService struct {
	Name     valueObject.ServiceName      `json:"name"`
	Nature   valueObject.ServiceNature    `json:"nature"`
	Type     valueObject.ServiceType      `json:"type"`
	Versions []valueObject.ServiceVersion `json:"versions"`
}

func NewInstallableService(
	name valueObject.ServiceName,
	serviceNature valueObject.ServiceNature,
	serviceType valueObject.ServiceType,
	versions []valueObject.ServiceVersion,
) InstallableService {
	return InstallableService{
		Name:     name,
		Nature:   serviceNature,
		Type:     serviceType,
		Versions: versions,
	}
}
