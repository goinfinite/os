package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func DeleteActivityRecord(
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteActivityRecord,
) error {
	err := activityRecordCmdRepo.Delete(deleteDto)
	if err != nil {
		slog.Error("DeleteActivityRecordError", slog.Any("err", err))
		return errors.New("DeleteActivityRecordInfraError")
	}

	return nil
}
