package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func CreateScheduledTask(
	scheduledTaskCmdRepo repository.ScheduledTaskCmdRepo,
	dto dto.CreateScheduledTask,
) error {
	err := scheduledTaskCmdRepo.Create(dto)
	if err != nil {
		slog.Error("CreateScheduledTaskError", slog.Any("error", err))
		return errors.New("CreateScheduledTaskInfraError")
	}

	return nil
}
