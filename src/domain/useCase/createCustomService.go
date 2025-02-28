package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func createFirstMapping(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	serviceName valueObject.ServiceName,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
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
		slog.Error("ReadPrimaryVhostMappingsError", slog.Any("error", err))
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
		primaryVhost.Hostname, mappingPath, matchPattern, targetType, &targetValue, nil,
		operatorAccountId, operatorIpAddress,
	)

	_, err = mappingCmdRepo.Create(createMappingDto)
	if err != nil {
		slog.Error("CreateServiceMappingError", slog.Any("error", err))
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
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateCustomService,
) error {
	readFirstInstalledRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &createDto.Name,
	}
	_, err := servicesQueryRepo.ReadFirstInstalledItem(readFirstInstalledRequestDto)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	if createDto.Version == nil {
		defaultVersion, _ := valueObject.NewServiceVersion("latest")
		createDto.Version = &defaultVersion
	}

	if createDto.Type == valueObject.SystemServiceType {
		return errors.New("SystemServiceCannotBeCreated")
	}

	err = servicesCmdRepo.CreateCustom(createDto)
	if err != nil {
		slog.Error("CreateCustomServiceError", slog.Any("error", err))
		return errors.New("CreateCustomServiceInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateCustomService(createDto)

	if createDto.AutoCreateMapping != nil && !*createDto.AutoCreateMapping {
		return nil
	}

	serviceTypeStr := createDto.Type.String()
	if serviceTypeStr != "runtime" && serviceTypeStr != "application" {
		return nil
	}

	return createFirstMapping(
		vhostQueryRepo, mappingQueryRepo, mappingCmdRepo, createDto.Name,
		createDto.OperatorAccountId, createDto.OperatorIpAddress,
	)
}
