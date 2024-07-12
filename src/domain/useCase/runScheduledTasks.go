package useCase

import (
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

const ScheduledTasksRunIntervalSecs uint = 120

var scheduledTasksDefaultTimeoutSecs uint = 300

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
			if nowUnixTime.Read() < pendingTask.RunAt.Read() {
				continue
			}
		}

		if pendingTask.TimeoutSecs == nil {
			pendingTask.TimeoutSecs = &scheduledTasksDefaultTimeoutSecs
		}

		err = scheduledTaskCmdRepo.Run(pendingTask)
		if err != nil {
			log.Printf("(%d) RunScheduledTaskError: %s", pendingTask.Id.Read(), err)
			continue
		}
	}
}
