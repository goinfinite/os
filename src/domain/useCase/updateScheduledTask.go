package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func UpdateScheduledTask(
	scheduledTaskQueryRepo repository.ScheduledTaskQueryRepo,
	scheduledTaskCmdRepo repository.ScheduledTaskCmdRepo,
	updateDto dto.UpdateScheduledTask,
) error {
	readDto := dto.ReadScheduledTasksRequest{
		Pagination: ScheduledTasksDefaultPagination,
		TaskId:     &updateDto.TaskId,
	}

	responseDto, err := scheduledTaskQueryRepo.Read(readDto)
	if err != nil {
		return errors.New("ReadScheduledTaskInfraError")
	}

	if len(responseDto.Tasks) == 0 {
		return errors.New("ScheduledTaskNotFound")
	}

	taskEntity := responseDto.Tasks[0]

	if taskEntity.Status == *updateDto.Status {
		slog.Debug("IgnoringScheduledTaskUpdateStatusNotChanged")
		return nil
	}

	if taskEntity.Status.String() == "running" {
		return errors.New("CannotUpdateRunningTask")
	}

	err = scheduledTaskCmdRepo.Update(updateDto)
	if err != nil {
		slog.Error("UpdateScheduledTaskInfraError", slog.Any("error", err))
		return errors.New("UpdateScheduledTaskInfraError")
	}

	return nil
}
