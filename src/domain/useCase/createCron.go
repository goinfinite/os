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
		slog.Error("CreateCronError", slog.Any("err", err))
		return errors.New("CreateCronInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateCron(createDto, cronId)

	cronCmdLimitStr := len(createDto.Command.String())
	if cronCmdLimitStr > 75 {
		cronCmdLimitStr = 75
	}
	cronCmdShortVersion := createDto.Command.String()[:cronCmdLimitStr]
	cronLine := createDto.Schedule.String() + " " + cronCmdShortVersion

	slog.Info("CronCreated", slog.String("cron", cronLine))

	return nil
}
