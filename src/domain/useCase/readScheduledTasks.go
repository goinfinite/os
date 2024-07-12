package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadScheduledTasks(
	scheduledTaskQueryRepo repository.ScheduledTaskQueryRepo,
) ([]entity.ScheduledTask, error) {
	scheduledTasks, err := scheduledTaskQueryRepo.Read()
	if err != nil {
		log.Printf("GetTasksError: %s", err)
		return nil, errors.New("GetTasksInfraError")
	}

	return scheduledTasks, nil
}
