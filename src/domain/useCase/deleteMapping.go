package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteMapping(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteMapping,
) error {
	_, err := mappingQueryRepo.ReadFirst(dto.ReadMappingsRequest{
		MappingId: &deleteDto.MappingId,
	})
	if err != nil {
		return errors.New("MappingNotFound")
	}

	err = mappingCmdRepo.Delete(deleteDto.MappingId)
	if err != nil {
		slog.Error("DeleteMappingError", slog.String("err", err.Error()))
		return errors.New("DeleteMappingInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteMapping(deleteDto)

	return nil
}
