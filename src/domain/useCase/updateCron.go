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
	readFirstRequestDto := dto.ReadCronsRequest{
		CronId: &updateDto.Id,
	}
	_, err := cronQueryRepo.ReadFirst(readFirstRequestDto)
	if err != nil {
		return errors.New("CronNotFound")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).UpdateCron(updateDto)

	err = cronCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateCronError", slog.String("err", err.Error()))
		return errors.New("UpdateCronInfraError")
	}

	return nil
}
