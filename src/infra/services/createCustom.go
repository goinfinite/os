package servicesInfra

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func CreateCustom(
	createDto dto.CreateCustomService,
) error {
	svcVersion := valueObject.NewServiceVersionPanic("latest")
	if createDto.Version != nil {
		svcVersion = *createDto.Version
	}

	return SupervisordFacade{}.CreateConf(
		createDto.Name,
		valueObject.NewServiceNaturePanic("custom"),
		createDto.Type,
		svcVersion,
		createDto.Command,
		nil,
		createDto.PortBindings,
		nil,
	)
}
