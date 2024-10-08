package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func UpdateScheduledTask(
	scheduledTaskQueryRepo repository.ScheduledTaskQueryRepo,
	scheduledTaskCmdRepo repository.ScheduledTaskCmdRepo,
	updateDto dto.UpdateScheduledTask,
) error {
	taskEntity, err := scheduledTaskQueryRepo.ReadById(updateDto.Id)
	if err != nil {
		return errors.New("ScheduledTaskNotFound")
	}

	if taskEntity.Status == *updateDto.Status {
		return nil
	}

	if taskEntity.Status.String() == "running" {
		return errors.New("CannotUpdateRunningTask")
	}

	err = scheduledTaskCmdRepo.Update(updateDto)
	if err != nil {
		log.Printf("UpdateScheduledTaskError: %s", err)
		return errors.New("UpdateScheduledTaskInfraError")
	}

	log.Printf("ScheduledTaskId '%v' updated.", updateDto.Id)

	return nil
}
