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
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	addDto dto.CreateInstallableService,
) error {
	_, err := servicesQueryRepo.GetByName(addDto.Name)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	installableSvcs, err := servicesQueryRepo.GetInstallables()
	if err != nil {
		log.Printf("GetInstallableServicesError: %s", err.Error())
		return errors.New("GetInstallableServicesInfraError")
	}

	dtoServiceNameStr := addDto.Name.String()
	isNatureMulti := false
	for _, installableSvc := range installableSvcs {
		if installableSvc.Name.String() != dtoServiceNameStr {
			continue
		}

		isNatureMulti = installableSvc.Nature.String() == "multi"
		break
	}

	if isNatureMulti {
		newSvcName, err := servicesQueryRepo.GetMultiServiceName(addDto.Name, addDto.StartupFile)
		if err != nil {
			log.Printf("GetMultiServiceNameError: %s", err.Error())
			return errors.New("GetMultiServiceNameInfraError")
		}

		addDto.Name = newSvcName
	}

	err = servicesCmdRepo.AddInstallable(addDto)
	if err != nil {
		log.Printf("CreateInstallableServiceError: %v", err)
		return errors.New("CreateInstallableServiceInfraError")
	}

	vhostsWithMappings, err := vhostQueryRepo.GetWithMappings()
	if err != nil {
		log.Printf("GetVhostsWithMappingError: %s", err.Error())
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

	err = vhostCmdRepo.CreateMapping(serviceMapping)
	if err != nil {
		log.Printf("AddServiceMappingError: %s", err.Error())
		return errors.New("AddServiceMappingInfraError")
	}

	return nil
}
