package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func UpdateService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	updateDto dto.UpdateService,
) error {
	shouldIncludeMetrics := false
	readInstalledDto := dto.ReadInstalledServicesItemsRequest{
		ServiceName:          &updateDto.Name,
		ShouldIncludeMetrics: &shouldIncludeMetrics,
	}
	serviceEntity, err := servicesQueryRepo.ReadOneInstalledItem(readInstalledDto)
	if err != nil {
		return err
	}

	isSoloService := serviceEntity.Nature.String() == "solo"
	isSystemService := serviceEntity.Type.String() == "system"
	shouldUpdateStatus := updateDto.Status != nil
	if (isSoloService || isSystemService) && !shouldUpdateStatus {
		return errors.New("OnlyStatusUpdateAllowed")
	}

	shouldDelete := shouldUpdateStatus && updateDto.Status.String() == "uninstalled"
	if shouldDelete {
		return DeleteService(
			servicesQueryRepo, servicesCmdRepo, mappingCmdRepo, updateDto.Name,
		)
	}

	err = servicesCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateServiceError", slog.Any("error", err))
		return errors.New("UpdateServiceInfraError")
	}

	if len(updateDto.PortBindings) == 0 {
		return nil
	}

	err = mappingCmdRepo.RecreateByServiceName(updateDto.Name)
	if err != nil {
		slog.Error("RecreateMappingError", slog.Any("error", err))
	}

	return nil
}
