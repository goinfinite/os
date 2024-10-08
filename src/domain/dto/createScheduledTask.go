package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateScheduledTask struct {
	Name        valueObject.ScheduledTaskName  `json:"name"`
	Command     valueObject.UnixCommand        `json:"command"`
	Tags        []valueObject.ScheduledTaskTag `json:"tags"`
	TimeoutSecs *uint                          `json:"timeoutSecs"`
	RunAt       *valueObject.UnixTime          `json:"runAt"`
}

func NewCreateScheduledTask(
	name valueObject.ScheduledTaskName,
	command valueObject.UnixCommand,
	tags []valueObject.ScheduledTaskTag,
	timeoutSecs *uint,
	runAt *valueObject.UnixTime,
) CreateScheduledTask {
	return CreateScheduledTask{
		Name:        name,
		Command:     command,
		Tags:        tags,
		TimeoutSecs: timeoutSecs,
		RunAt:       runAt,
	}
}
