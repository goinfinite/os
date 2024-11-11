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
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateInstallableService,
) error {
	_, err := servicesQueryRepo.ReadByName(createDto.Name)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	installedServiceName, err := servicesCmdRepo.CreateInstallable(createDto)
	if err != nil {
		slog.Error("CreateInstallableServiceError", slog.Any("error", err))
		return errors.New("CreateInstallableServiceInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateInstallableService(createDto)

	serviceEntity, err := servicesQueryRepo.ReadByName(installedServiceName)
	if err != nil {
		slog.Error("GetServiceByNameError", slog.Any("error", err))
		return errors.New("GetServiceByNameInfraError")
	}

	if createDto.AutoCreateMapping != nil && !*createDto.AutoCreateMapping {
		return nil
	}

	serviceTypeStr := serviceEntity.Type.String()
	if serviceTypeStr != "runtime" && serviceTypeStr != "application" {
		return nil
	}

	return createFirstMapping(
		vhostQueryRepo, mappingQueryRepo, mappingCmdRepo, installedServiceName,
	)
}
