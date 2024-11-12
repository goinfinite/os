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
	mapping, err := mappingQueryRepo.ReadById(deleteDto.MappingId)
	if err != nil {
		return errors.New("MappingNotFound")
	}

	err = mappingCmdRepo.Delete(deleteDto.MappingId)
	if err != nil {
		slog.Error("DeleteMappingError", slog.Any("err", err))
		return errors.New("DeleteMappingInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteMapping(deleteDto)

	slog.Info(
		"MappingDeleted", slog.String("path", mapping.Path.String()),
		slog.String("hostname", mapping.Hostname.String()),
	)

	return nil
}
