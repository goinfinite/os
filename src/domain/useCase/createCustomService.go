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
	)

	return serviceMapping, nil
}

func CreateCustomService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	addDto dto.CreateCustomService,
) error {
	_, err := servicesQueryRepo.GetByName(addDto.Name)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	err = servicesCmdRepo.AddCustom(addDto)
	if err != nil {
		log.Printf("CreateCustomServiceError: %v", err)
		return errors.New("CreateCustomServiceInfraError")
	}

	isRuntimeSvc := addDto.Type.String() == "runtime"
	isApplicationSvc := addDto.Type.String() == "application"
	if !isRuntimeSvc && !isApplicationSvc {
		return nil
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

	err = vhostCmdRepo.CreateMapping(serviceMapping)
	if err != nil {
		log.Printf("AddServiceMappingError: %s", err.Error())
		return errors.New("AddServiceMappingInfraError")
	}

	return nil
}
