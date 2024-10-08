package useCase

import (
	"log"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

const ScheduledTasksRunIntervalSecs uint = 120

func RunScheduledTasks(
	scheduledTaskQueryRepo repository.ScheduledTaskQueryRepo,
	scheduledTaskCmdRepo repository.ScheduledTaskCmdRepo,
) {
	pendingStatus, _ := valueObject.NewScheduledTaskStatus("pending")
	pendingTasks, err := scheduledTaskQueryRepo.ReadByStatus(pendingStatus)
	if err != nil {
		log.Printf("GetPendingScheduledTasksError: %s", err)
		return
	}

	if len(pendingTasks) == 0 {
		return
	}

	for _, pendingTask := range pendingTasks {
		if pendingTask.RunAt != nil {
			nowUnixTime := valueObject.NewUnixTimeNow()
			if nowUnixTime.Int64() < pendingTask.RunAt.Int64() {
				continue
			}
		}

		err = scheduledTaskCmdRepo.Run(pendingTask)
		if err != nil {
			log.Printf("(%d) RunScheduledTaskError: %s", pendingTask.Id, err)
			continue
		}
	}
}
