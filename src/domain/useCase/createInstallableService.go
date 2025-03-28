package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CreateInstallableService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateInstallableService,
) error {
	_, err := servicesQueryRepo.ReadFirstInstalledItem(dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &createDto.Name,
	})
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	installedServiceName, err := servicesCmdRepo.CreateInstallable(createDto)
	if err != nil {
		slog.Error("CreateInstallableServiceError", slog.String("err", err.Error()))
		return errors.New("CreateInstallableServiceInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateInstallableService(createDto)

	serviceEntity, err := servicesQueryRepo.ReadFirstInstalledItem(dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &installedServiceName,
	})
	if err != nil {
		slog.Error("ReadServiceEntityError", slog.String("err", err.Error()))
		return errors.New("ReadServiceEntityInfraError")
	}

	if createDto.AutoCreateMapping != nil && !*createDto.AutoCreateMapping {
		return nil
	}

	if len(serviceEntity.PortBindings) == 0 {
		slog.Debug("AutoCreateMappingSkipped", slog.String("reason", "PortBindingsIsEmpty"))
		return nil
	}

	return CreateServiceAutoMapping(
		vhostQueryRepo, mappingCmdRepo, installedServiceName, createDto.MappingHostname,
		createDto.MappingPath, createDto.OperatorAccountId, createDto.OperatorIpAddress,
	)
}
