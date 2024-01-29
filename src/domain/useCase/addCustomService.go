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
) (dto.AddMapping, error) {
	var serviceMapping dto.AddMapping

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

	serviceMapping = dto.NewAddMapping(
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

func AddCustomService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	addDto dto.AddCustomService,
) error {
	_, err := servicesQueryRepo.GetByName(addDto.Name)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	err = servicesCmdRepo.AddCustom(addDto)
	if err != nil {
		log.Printf("AddCustomServiceError: %v", err)
		return errors.New("AddCustomServiceInfraError")
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

	err = vhostCmdRepo.AddMapping(serviceMapping)
	if err != nil {
		log.Printf("AddServiceMappingError: %s", err.Error())
		return errors.New("AddServiceMappingInfraError")
	}

	return nil
}
