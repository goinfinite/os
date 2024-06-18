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
	mappingCmdRepo repository.MappingCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
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

	serviceEntity, err := servicesQueryRepo.GetByName(createDto.Name)
	if err != nil {
		log.Printf("GetServiceByNameError: %s", err.Error())
		return errors.New("GetServiceByNameInfraError")
	}

	isRuntimeSvc := serviceEntity.Type.String() == "runtime"
	isApplicationSvc := serviceEntity.Type.String() == "application"
	if !isRuntimeSvc && !isApplicationSvc {
		return nil
	}

	vhosts, err := vhostQueryRepo.Read()
	if err != nil {
		return errors.New("VhostsNotFound")
	}

	primaryVhost := vhosts[0]
	primaryVhostMappings, err := mappingQueryRepo.ReadByHostname(
		primaryVhost.Hostname,
	)
	if err != nil {
		log.Printf("ReadPrimaryVhostMappingsError: %s", err.Error())
		return errors.New("ReadPrimaryVhostMappingsInfraError")
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

	_, err = mappingCmdRepo.Create(serviceMapping)
	if err != nil {
		log.Printf("CreateServiceMappingError: %s", err.Error())
		return errors.New("CreateServiceMappingInfraError")
	}

	return nil
}
