package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateScheduledTask struct {
	Name        valueObject.ScheduledTaskName  `json:"name"`
	Command     valueObject.UnixCommand        `json:"command"`
	Tags        []valueObject.ScheduledTaskTag `json:"tags"`
	TimeoutSecs *uint16                        `json:"timeoutSecs,omitempty"`
	RunAt       *valueObject.UnixTime          `json:"runAt,omitempty"`
}

func NewCreateScheduledTask(
	name valueObject.ScheduledTaskName,
	command valueObject.UnixCommand,
	tags []valueObject.ScheduledTaskTag,
	timeoutSecs *uint16,
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
