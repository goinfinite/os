package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CreateCron(
	cronCmdRepo repository.CronCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateCron,
) error {
	cronId, err := cronCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateCronError", slog.String("err", err.Error()))
		return errors.New("CreateCronInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateCron(createDto, cronId)

	return nil
}
