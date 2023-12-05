package entity

import "github.com/speedianet/os/src/domain/valueObject"

type InstallableService struct {
	Name     valueObject.ServiceName      `json:"name"`
	Type     valueObject.ServiceType      `json:"type"`
	Versions []valueObject.ServiceVersion `json:"versions"`
}

func NewInstallableService(
	name valueObject.ServiceName,
	serviceType valueObject.ServiceType,
	versions []valueObject.ServiceVersion,
) InstallableService {
	return InstallableService{
		Name:     name,
		Type:     serviceType,
		Versions: versions,
	}
}
