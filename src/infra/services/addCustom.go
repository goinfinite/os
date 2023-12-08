package servicesInfra

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func AddCustom(
	addDto dto.AddCustomService,
) error {
	return SupervisordFacade{}.AddConf(
		addDto.Name,
		valueObject.NewServiceNaturePanic("custom"),
		addDto.Type,
		addDto.Command,
		addDto.Ports,
	)
}
