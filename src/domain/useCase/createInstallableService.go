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
	readFirstInstalledRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &createDto.Name,
	}
	_, err := servicesQueryRepo.ReadFirstInstalledItem(readFirstInstalledRequestDto)
	if err == nil {
		return errors.New("ServiceAlreadyInstalled")
	}

	installedServiceName, err := servicesCmdRepo.CreateInstallable(createDto)
	if err != nil {
		slog.Error("CreateInstallableServiceError", slog.Any("error", err))
		return errors.New("CreateInstallableServiceInfraError")
	}

	readFirstInstalledRequestDto.ServiceName = &installedServiceName
	serviceEntity, err := servicesQueryRepo.ReadFirstInstalledItem(
		readFirstInstalledRequestDto,
	)
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
