package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CreateInstallableService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	createDto dto.CreateInstallableService,
) error {
	_, err := servicesQueryRepo.GetByName(createDto.Name)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	installableSvcs, err := servicesQueryRepo.GetInstallables()
	if err != nil {
		log.Printf("GetInstallableServicesError: %s", err.Error())
		return errors.New("GetInstallableServicesInfraError")
	}

	dtoServiceNameStr := createDto.Name.String()
	isNatureMulti := false
	for _, installableSvc := range installableSvcs {
		if installableSvc.Name.String() != dtoServiceNameStr {
			continue
		}

		isNatureMulti = installableSvc.Nature.String() == "multi"
		break
	}

	if isNatureMulti {
		newSvcName, err := servicesQueryRepo.GetMultiServiceName(createDto.Name, createDto.StartupFile)
		if err != nil {
			log.Printf("GetMultiServiceNameError: %s", err.Error())
			return errors.New("GetMultiServiceNameInfraError")
		}

		createDto.Name = newSvcName
	}

	err = servicesCmdRepo.CreateInstallable(createDto)
	if err != nil {
		log.Printf("CreateInstallableServiceError: %v", err)
		return errors.New("CreateInstallableServiceInfraError")
	}

	vhosts, err := vhostQueryRepo.Get()
	if err != nil {
		return errors.New("VhostsNotFound")
	}

	primaryVhost := vhosts[0]
	primaryVhostMappings, err := mappingQueryRepo.GetByHostname(
		primaryVhost.Hostname,
	)
	if err != nil {
		log.Printf("GetPrimaryVhostMappingsError: %s", err.Error())
		return errors.New("GetPrimaryVhostMappingsInfraError")
	}
	shouldCreateFirstMapping := len(primaryVhostMappings) == 0 && createDto.AutoCreateMapping
	if !shouldCreateFirstMapping {
		return nil
	}

	serviceMapping, err := serviceMappingFactory(
		primaryVhost.Hostname,
		createDto.Name,
	)
	if err != nil {
		log.Printf("CreateServiceMappingError: %s", err.Error())
		return errors.New("CreateServiceMappingError")
	}

	err = vhostCmdRepo.CreateMapping(serviceMapping)
	if err != nil {
		log.Printf("CreateServiceMappingError: %s", err.Error())
		return errors.New("CreateServiceMappingInfraError")
	}

	return nil
}
