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
	createDto dto.CreateInstallableService,
) error {
	readInstalledDto := dto.ReadInstalledServicesItemsRequest{
		Name:                 &createDto.Name,
		ShouldIncludeMetrics: false,
	}
	_, err := servicesQueryRepo.ReadUniqueInstalledItem(readInstalledDto)
	if err != nil {
		return err
	}

	installedServiceName, err := servicesCmdRepo.CreateInstallable(createDto)
	if err != nil {
		slog.Error("CreateInstallableServiceError", slog.Any("error", err))
		return errors.New("CreateInstallableServiceInfraError")
	}

	readInstalledDto.Name = &installedServiceName
	serviceEntity, err := servicesQueryRepo.ReadUniqueInstalledItem(readInstalledDto)
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
