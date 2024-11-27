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
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateService,
) error {
	readFirstInstalledRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &updateDto.Name,
	}
	serviceEntity, err := servicesQueryRepo.ReadFirstInstalledItem(
		readFirstInstalledRequestDto,
	)
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
		deleteDto := dto.NewDeleteService(
			updateDto.Name, updateDto.OperatorAccountId, updateDto.OperatorIpAddress,
		)
		return DeleteService(
			servicesQueryRepo, servicesCmdRepo, mappingCmdRepo, activityRecordCmdRepo,
			deleteDto,
		)
	}

	err = servicesCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateServiceError", slog.Any("error", err))
		return errors.New("UpdateServiceInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateService(updateDto)

	if len(updateDto.PortBindings) == 0 {
		return nil
	}

	err = mappingCmdRepo.RecreateByServiceName(
		updateDto.Name, updateDto.OperatorAccountId, updateDto.OperatorIpAddress,
	)
	if err != nil {
		slog.Error("RecreateMappingError", slog.Any("error", err))
	}

	return nil
}
