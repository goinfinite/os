package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func UpdateMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateMapping,
) error {
	_, err := mappingQueryRepo.ReadFirst(dto.ReadMappingsRequest{
		MappingId: &updateDto.Id,
	})
	if err != nil {
		slog.Debug("ReadMappingError", slog.String("err", err.Error()))
		return errors.New("MappingNotFound")
	}

	if updateDto.TargetType != nil &&
		*updateDto.TargetType == valueObject.MappingTargetTypeResponseCode {
		// TODO: When the "truncatable" feature is implemented, change this to
		// populate the truncatable slice instead.
		updateDto.TargetValue = nil
	}

	err = mappingCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateMappingError", slog.String("err", err.Error()))
		return errors.New("UpdateMappingInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateMapping(updateDto)

	return nil
}
