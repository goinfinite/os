package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateScheduledTask struct {
	Id     valueObject.ScheduledTaskId      `json:"id"`
	Status *valueObject.ScheduledTaskStatus `json:"status"`
	RunAt  *valueObject.UnixTime            `json:"runAt"`
}

func NewUpdateScheduledTask(
	id valueObject.ScheduledTaskId,
	status *valueObject.ScheduledTaskStatus,
	runAt *valueObject.UnixTime,
) UpdateScheduledTask {
	return UpdateScheduledTask{
		Id:     id,
		Status: status,
		RunAt:  runAt,
	}
}
