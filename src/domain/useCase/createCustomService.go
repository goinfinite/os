package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func createFirstMapping(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	serviceName valueObject.ServiceName,
) error {
	vhosts, err := vhostQueryRepo.Read()
	if err != nil {
		return errors.New("VhostsNotFound")
	}

	primaryVhost := vhosts[0]
	primaryVhostMappings, err := mappingQueryRepo.ReadByHostname(
		primaryVhost.Hostname,
	)
	if err != nil {
		slog.Error("ReadPrimaryVhostMappingsError", slog.Any("err", err))
		return errors.New("ReadPrimaryVhostMappingsInfraError")
	}
	if len(primaryVhostMappings) != 0 {
		return nil
	}

	mappingPath, _ := valueObject.NewMappingPath("/")
	matchPattern, _ := valueObject.NewMappingMatchPattern("begins-with")
	targetType, _ := valueObject.NewMappingTargetType("service")
	targetValue, _ := valueObject.NewMappingTargetValue(serviceName.String(), targetType)

	createMappingDto := dto.NewCreateMapping(
		primaryVhost.Hostname,
		mappingPath,
		matchPattern,
		targetType,
		&targetValue,
		nil,
	)

	_, err = mappingCmdRepo.Create(createMappingDto)
	if err != nil {
		slog.Error("CreateServiceMappingError", slog.Any("err", err))
		return errors.New("CreateServiceMappingInfraError")
	}

	return nil
}

func CreateCustomService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	createDto dto.CreateCustomService,
) error {
	_, err := servicesQueryRepo.ReadByName(createDto.Name)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	if createDto.Version == nil {
		defaultVersion, _ := valueObject.NewServiceVersion("latest")
		createDto.Version = &defaultVersion
	}

	err = servicesCmdRepo.CreateCustom(createDto)
	if err != nil {
		slog.Error("CreateCustomServiceError", slog.Any("err", err))
		return errors.New("CreateCustomServiceInfraError")
	}

	if createDto.AutoCreateMapping != nil && !*createDto.AutoCreateMapping {
		return nil
	}

	serviceTypeStr := createDto.Type.String()
	if serviceTypeStr != "runtime" && serviceTypeStr != "application" {
		return nil
	}

	return createFirstMapping(
		vhostQueryRepo, mappingQueryRepo, mappingCmdRepo, createDto.Name,
	)
}
