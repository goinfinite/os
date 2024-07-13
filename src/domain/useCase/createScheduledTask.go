package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

var ScheduledTasksDefaultTimeoutSecs uint = 300

func CreateScheduledTask(
	scheduledTaskCmdRepo repository.ScheduledTaskCmdRepo,
	dto dto.CreateScheduledTask,
) error {
	if dto.TimeoutSecs == nil {
		dto.TimeoutSecs = &ScheduledTasksDefaultTimeoutSecs
	}

	err := scheduledTaskCmdRepo.Create(dto)
	if err != nil {
		log.Printf("CreateScheduledTaskError: %v", err)
		return errors.New("CreateScheduledTaskInfraError")
	}

	return nil
}
