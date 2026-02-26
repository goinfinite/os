package useCase

import (
	"errors"
	"log/slog"

	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func DeleteActivityRecord(
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	deleteDto tkDto.DeleteActivityRecord,
) error {
	err := activityRecordCmdRepo.Delete(deleteDto)
	if err != nil {
		slog.Error("DeleteActivityRecordError", slog.String("err", err.Error()))
		return errors.New("DeleteActivityRecordInfraError")
	}

	return nil
}
