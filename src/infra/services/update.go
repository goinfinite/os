package servicesInfra

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
)

func Update(
	serviceEntity entity.Service,
	updateDto dto.UpdateService,
) error {
	err := SupervisordFacade{}.RemoveConf(updateDto.Name)
	if err != nil {
		return errors.New("RemoveServiceConfError")
	}

	svcType := serviceEntity.Type
	if updateDto.Type != nil {
		svcType = *updateDto.Type
	}

	svcVersion := serviceEntity.Version
	if updateDto.Version != nil {
		svcVersion = *updateDto.Version
	}

	svcCommand := serviceEntity.Command
	if updateDto.Command != nil {
		svcCommand = *updateDto.Command
	}

	svcStartupFile := serviceEntity.StartupFile
	if updateDto.StartupFile != nil {
		svcStartupFile = updateDto.StartupFile
	}

	svcPorts := serviceEntity.Ports
	if len(updateDto.Ports) > 0 {
		svcPorts = updateDto.Ports
	}

	err = SupervisordFacade{}.AddConf(
		serviceEntity.Name,
		serviceEntity.Nature,
		svcType,
		svcVersion,
		svcCommand,
		svcStartupFile,
		svcPorts,
	)
	if err != nil {
		return errors.New("ReAddServiceConfError")
	}

	return nil
}
