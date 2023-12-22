package servicesInfra

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func AddCustom(
	addDto dto.AddCustomService,
) error {
	svcVersion := valueObject.NewServiceVersionPanic("latest")
	if addDto.Version != nil {
		svcVersion = *addDto.Version
	}

	return SupervisordFacade{}.AddConf(
		addDto.Name,
		valueObject.NewServiceNaturePanic("custom"),
		addDto.Type,
		svcVersion,
		addDto.Command,
		nil,
		addDto.PortBindings,
	)
}
