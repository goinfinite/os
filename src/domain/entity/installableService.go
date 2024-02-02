package entity

import "github.com/speedianet/os/src/domain/valueObject"

type InstallableService struct {
	Name     valueObject.ServiceName      `json:"name"`
	Type     valueObject.ServiceType      `json:"type"`
	Nature   valueObject.ServiceNature    `json:"nature"`
	Versions []valueObject.ServiceVersion `json:"versions"`
}

func NewInstallableService(
	name valueObject.ServiceName,
	serviceType valueObject.ServiceType,
	serviceNature valueObject.ServiceNature,
	versions []valueObject.ServiceVersion,
) InstallableService {
	return InstallableService{
		Name:     name,
		Type:     serviceType,
		Nature:   serviceNature,
		Versions: versions,
	}
}
