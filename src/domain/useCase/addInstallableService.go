package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func AddInstallableService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	addDto dto.AddInstallableService,
) error {
	_, err := servicesQueryRepo.GetByName(addDto.Name)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	err = servicesCmdRepo.AddInstallable(addDto)
	if err != nil {
		log.Printf("AddInstallableServiceError: %v", err)
		return errors.New("AddInstallableServiceInfraError")
	}

	vhostsWithMappings, err := vhostQueryRepo.GetWithMappings()
	if err != nil {
		return errors.New("GetVhostsWithMappingsInfraError")
	}

	if len(vhostsWithMappings) == 0 {
		return errors.New("VhostsNotFound")
	}

	primaryVhostWithMapping := vhostsWithMappings[0]
	shouldCreateFirstMapping := len(primaryVhostWithMapping.Mappings) == 0 && addDto.AutoCreateMapping
	if !shouldCreateFirstMapping {
		return nil
	}

	serviceMapping, err := serviceMappingFactory(
		primaryVhostWithMapping.Hostname,
		addDto.Name,
	)
	if err != nil {
		log.Printf("AddServiceMappingError: %s", err.Error())
		return errors.New("AddServiceMappingError")
	}

	err = vhostCmdRepo.AddMapping(serviceMapping)
	if err != nil {
		log.Printf("AddServiceMappingError: %s", err.Error())
		return errors.New("AddServiceMappingInfraError")
	}

	return nil
}
