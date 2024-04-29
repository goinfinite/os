package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func serviceMappingFactory(
	primaryHostname valueObject.Fqdn,
	svcName valueObject.ServiceName,
) (dto.CreateMapping, error) {
	var serviceMapping dto.CreateMapping

	svcMappingPath, err := valueObject.NewMappingPath("/")
	if err != nil {
		return serviceMapping, err
	}

	svcMappingMatchPattern, err := valueObject.NewMappingMatchPattern("begins-with")
	if err != nil {
		return serviceMapping, err
	}

	svcMappingTargetType, err := valueObject.NewMappingTargetType("service")
	if err != nil {
		return serviceMapping, err
	}

	serviceMapping = dto.NewCreateMapping(
		primaryHostname,
		svcMappingPath,
		svcMappingMatchPattern,
		svcMappingTargetType,
		&svcName,
		nil,
		nil,
		nil,
	)

	return serviceMapping, nil
}

func CreateCustomService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	createDto dto.CreateCustomService,
) error {
	_, err := servicesQueryRepo.GetByName(createDto.Name)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	err = servicesCmdRepo.CreateCustom(createDto)
	if err != nil {
		log.Printf("CreateCustomServiceError: %v", err)
		return errors.New("CreateCustomServiceInfraError")
	}

	isRuntimeSvc := createDto.Type.String() == "runtime"
	isApplicationSvc := createDto.Type.String() == "application"
	if !isRuntimeSvc && !isApplicationSvc {
		return nil
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

	_, err = mappingCmdRepo.Create(serviceMapping)
	if err != nil {
		log.Printf("CreateServiceMappingError: %s", err.Error())
		return errors.New("CreateServiceMappingInfraError")
	}

	return nil
}
