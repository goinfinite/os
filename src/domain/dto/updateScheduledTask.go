package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateScheduledTask struct {
	TaskId valueObject.ScheduledTaskId      `json:"taskId"`
	Status *valueObject.ScheduledTaskStatus `json:"status,omitempty"`
	RunAt  *valueObject.UnixTime            `json:"runAt,omitempty"`
}

func NewUpdateScheduledTask(
	taskId valueObject.ScheduledTaskId,
	status *valueObject.ScheduledTaskStatus,
	runAt *valueObject.UnixTime,
) UpdateScheduledTask {
	return UpdateScheduledTask{
		TaskId: taskId,
		Status: status,
		RunAt:  runAt,
	}
}
