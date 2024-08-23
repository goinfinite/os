package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func UpdateService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	updateDto dto.UpdateService,
) error {
	serviceEntity, err := servicesQueryRepo.ReadByName(updateDto.Name)
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
		slog.Error("UpdateServiceError", slog.Any("err", err))
		return errors.New("UpdateServiceInfraError")
	}

	if len(updateDto.PortBindings) == 0 {
		return nil
	}

	err = mappingCmdRepo.RecreateByServiceName(updateDto.Name)
	if err != nil {
		slog.Error("RecreateMappingError", slog.Any("err", err))
	}

	return nil
}
