package servicesInfra

import (
	"errors"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
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
		newCommandStr := strings.Replace(
			svcCommand.String(),
			serviceEntity.StartupFile.String(),
			updateDto.StartupFile.String(),
			1,
		)
		svcCommand = valueObject.UnixCommand(newCommandStr)
		svcStartupFile = updateDto.StartupFile
	}

	svcPortBindings := serviceEntity.PortBindings
	if len(updateDto.PortBindings) > 0 {
		svcPortBindings = updateDto.PortBindings
	}

	err = SupervisordFacade{}.AddConf(
		serviceEntity.Name,
		serviceEntity.Nature,
		svcType,
		svcVersion,
		svcCommand,
		svcStartupFile,
		svcPortBindings,
		nil,
	)
	if err != nil {
		return errors.New("ReAddServiceConfError")
	}

	return nil
}
