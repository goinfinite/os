package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func UpdateCron(
	cronQueryRepo repository.CronQueryRepo,
	cronCmdRepo repository.CronCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateCron,
) error {
	_, err := cronQueryRepo.ReadById(updateDto.Id)
	if err != nil {
		slog.Error("CronNotFound", slog.Any("err", err))
		return errors.New("CronNotFound")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateCron(updateDto)

	err = cronCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateCronError", slog.Any("err", err))
		return errors.New("UpdateCronInfraError")
	}

	return nil
}
