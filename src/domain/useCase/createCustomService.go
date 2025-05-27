package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func CreateServiceAutoMapping(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	serviceName valueObject.ServiceName,
	mappingHostname *valueObject.Fqdn,
	mappingPath *valueObject.MappingPath,
	mappingUpgradeInsecureRequests *bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) error {
	if mappingHostname == nil {
		isPrimary := true
		primaryVirtualHost, err := vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
			IsPrimary: &isPrimary,
		})
		if err != nil {
			return errors.New("ReadPrimaryVirtualHostError: " + err.Error())
		}
		mappingHostname = &primaryVirtualHost.Hostname
	}

	withMappings := true
	vhostWithMappings, err := vhostQueryRepo.ReadFirstWithMappings(dto.ReadVirtualHostsRequest{
		Hostname:     mappingHostname,
		WithMappings: &withMappings,
	},
	)
	if err != nil {
		return errors.New("ReadFirstVirtualHostWithMappingsError: " + err.Error())
	}

	if mappingPath == nil {
		rootPath, _ := valueObject.NewMappingPath("/")
		mappingPath = &rootPath
	}
	for _, mappingEntity := range vhostWithMappings.Mappings {
		if mappingEntity.Path == *mappingPath {
			slog.Debug(
				"MappingAlreadyExists",
				slog.String("method", "CreateServiceAutoMapping"),
				slog.String("hostname", mappingHostname.String()),
				slog.String("path", mappingPath.String()),
			)
			return nil
		}
	}

	targetValue, err := valueObject.NewMappingTargetValue(
		serviceName.String(), valueObject.MappingTargetTypeService,
	)
	if err != nil {
		return errors.New("NewMappingTargetValueError: " + err.Error())
	}

	_, err = mappingCmdRepo.Create(dto.NewCreateMapping(
		*mappingHostname, *mappingPath, valueObject.MappingMatchPatternBeginsWith,
		valueObject.MappingTargetTypeService, &targetValue, nil,
		mappingUpgradeInsecureRequests, nil, operatorAccountId, operatorIpAddress,
	))
	if err != nil {
		return errors.New("CreateServiceMappingInfraError: " + err.Error())
	}

	return nil
}

func CreateCustomService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateCustomService,
) error {
	_, err := servicesQueryRepo.ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest{ServiceName: &createDto.Name},
	)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	if createDto.Version == nil {
		defaultVersion, _ := valueObject.NewServiceVersion("latest")
		createDto.Version = &defaultVersion
	}

	if createDto.Type == valueObject.ServiceTypeSystem {
		return errors.New("SystemServiceCannotBeCreated")
	}

	err = servicesCmdRepo.CreateCustom(createDto)
	if err != nil {
		slog.Error("CreateCustomServiceError", slog.String("err", err.Error()))
		return errors.New("CreateCustomServiceInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateCustomService(createDto)

	if createDto.AutoCreateMapping != nil && !*createDto.AutoCreateMapping {
		return nil
	}

	if len(createDto.PortBindings) == 0 {
		slog.Debug("AutoCreateMappingSkipped", slog.String("reason", "PortBindingsIsEmpty"))
		return nil
	}

	err = CreateServiceAutoMapping(
		vhostQueryRepo, mappingCmdRepo, createDto.Name, createDto.MappingHostname,
		createDto.MappingPath, createDto.MappingUpgradeInsecureRequests,
		createDto.OperatorAccountId, createDto.OperatorIpAddress,
	)
	if err != nil {
		slog.Error("AutoCreateMappingError", slog.String("err", err.Error()))
		return errors.New("AutoCreateMappingInfraError")
	}

	return nil
}
