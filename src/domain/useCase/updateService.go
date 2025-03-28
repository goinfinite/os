package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func UpdateService(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateService,
) error {
	serviceEntity, err := servicesQueryRepo.ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest{ServiceName: &updateDto.Name},
	)
	if err != nil {
		slog.Error("ReadServiceInfraEntityError", slog.String("err", err.Error()))
		return errors.New("ReadServiceEntityError")
	}

	isSoloService := serviceEntity.Nature == valueObject.ServiceNatureSolo
	isSystemService := serviceEntity.Type == valueObject.ServiceTypeSystem
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
			servicesQueryRepo, servicesCmdRepo, mappingQueryRepo, mappingCmdRepo,
			activityRecordCmdRepo, deleteDto,
		)
	}

	err = servicesCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateServiceError", slog.String("err", err.Error()))
		return errors.New("UpdateServiceInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateService(updateDto)

	if len(updateDto.PortBindings) == 0 {
		return nil
	}

	shouldRecreateAutoMappings := true
	return DeleteServiceMappings(
		mappingQueryRepo, mappingCmdRepo, updateDto.Name, shouldRecreateAutoMappings,
		updateDto.OperatorAccountId, updateDto.OperatorIpAddress,
	)
}
