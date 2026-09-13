package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func CreateScheduledTask(
	scheduledTaskCmdRepo repository.ScheduledTaskCmdRepo,
	dto dto.CreateScheduledTask,
) (valueObject.ScheduledTaskId, error) {
	taskId, err := scheduledTaskCmdRepo.Create(dto)
	if err != nil {
		slog.Error("CreateScheduledTaskError", slog.String("err", err.Error()))
		return 0, errors.New("CreateScheduledTaskInfraError")
	}

	return taskId, nil
}
