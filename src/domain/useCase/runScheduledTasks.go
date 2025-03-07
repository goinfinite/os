package useCase

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

const ScheduledTasksRunIntervalSecs uint8 = 90

var scheduledTasksDefaultTimeoutSecs uint16 = 300

func RunScheduledTasks(
	scheduledTaskQueryRepo repository.ScheduledTaskQueryRepo,
	scheduledTaskCmdRepo repository.ScheduledTaskCmdRepo,
) {
	pendingStatus, _ := valueObject.NewScheduledTaskStatus("pending")
	readDto := dto.ReadScheduledTasksRequest{
		Pagination: ScheduledTasksDefaultPagination,
		TaskStatus: &pendingStatus,
	}

	responseDto, err := scheduledTaskQueryRepo.Read(readDto)
	if err != nil {
		slog.Error("ReadPendingScheduledTasksError", slog.String("err", err.Error()))
		return
	}

	if len(responseDto.Tasks) == 0 {
		return
	}

	for _, pendingTask := range responseDto.Tasks {
		if pendingTask.RunAt != nil {
			nowUnixTime := valueObject.NewUnixTimeNow()
			if nowUnixTime.Int64() < pendingTask.RunAt.Int64() {
				continue
			}
		}

		if pendingTask.TimeoutSecs == nil {
			pendingTask.TimeoutSecs = &scheduledTasksDefaultTimeoutSecs
		}

		err = scheduledTaskCmdRepo.Run(pendingTask)
		if err != nil {
			slog.Error(
				"RunScheduledTaskError",
				slog.Uint64("scheduledTaskId", pendingTask.Id.Uint64()),
				slog.String("err", err.Error()),
			)
			continue
		}
	}
}
