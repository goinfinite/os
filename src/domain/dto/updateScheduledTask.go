package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdateScheduledTask struct {
	TaskId valueObject.ScheduledTaskId      `json:"taskId"`
	Status *valueObject.ScheduledTaskStatus `json:"status,omitempty"`
	RunAt  *tkValueObject.UnixTime          `json:"runAt,omitempty"`
}

func NewUpdateScheduledTask(
	taskId valueObject.ScheduledTaskId,
	status *valueObject.ScheduledTaskStatus,
	runAt *tkValueObject.UnixTime,
) UpdateScheduledTask {
	return UpdateScheduledTask{
		TaskId: taskId,
		Status: status,
		RunAt:  runAt,
	}
}
