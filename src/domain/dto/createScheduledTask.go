package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateScheduledTask struct {
	Name        valueObject.ScheduledTaskName  `json:"name"`
	Command     tkValueObject.UnixCommand      `json:"command"`
	Tags        []valueObject.ScheduledTaskTag `json:"tags"`
	TimeoutSecs *uint16                        `json:"timeoutSecs,omitempty"`
	RunAt       *tkValueObject.UnixTime        `json:"runAt,omitempty"`
}

func NewCreateScheduledTask(
	name valueObject.ScheduledTaskName,
	command tkValueObject.UnixCommand,
	tags []valueObject.ScheduledTaskTag,
	timeoutSecs *uint16,
	runAt *tkValueObject.UnixTime,
) CreateScheduledTask {
	return CreateScheduledTask{
		Name:        name,
		Command:     command,
		Tags:        tags,
		TimeoutSecs: timeoutSecs,
		RunAt:       runAt,
	}
}
