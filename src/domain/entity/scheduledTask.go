package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ScheduledTask struct {
	Id          valueObject.ScheduledTaskId      `json:"id"`
	Name        valueObject.ScheduledTaskName    `json:"name"`
	Status      valueObject.ScheduledTaskStatus  `json:"status"`
	Command     tkValueObject.UnixCommand        `json:"command"`
	Tags        []valueObject.ScheduledTaskTag   `json:"tags"`
	TimeoutSecs *uint16                          `json:"timeoutSecs"`
	RunAt       *tkValueObject.UnixTime          `json:"runAt"`
	Output      *valueObject.ScheduledTaskOutput `json:"output"`
	Error       *valueObject.ScheduledTaskOutput `json:"err"`
	StartedAt   *tkValueObject.UnixTime          `json:"startedAt"`
	FinishedAt  *tkValueObject.UnixTime          `json:"finishedAt"`
	ElapsedSecs *uint32                          `json:"elapsedSecs"`
	CreatedAt   tkValueObject.UnixTime           `json:"createdAt"`
	UpdatedAt   tkValueObject.UnixTime           `json:"updatedAt"`
}

func NewScheduledTask(
	id valueObject.ScheduledTaskId,
	name valueObject.ScheduledTaskName,
	status valueObject.ScheduledTaskStatus,
	command tkValueObject.UnixCommand,
	tags []valueObject.ScheduledTaskTag,
	timeoutSecs *uint16,
	runAt *tkValueObject.UnixTime,
	output, err *valueObject.ScheduledTaskOutput,
	startedAt, finishedAt *tkValueObject.UnixTime,
	elapsedSecs *uint32,
	createdAt, updatedAt tkValueObject.UnixTime,
) ScheduledTask {
	return ScheduledTask{
		Id:          id,
		Name:        name,
		Status:      status,
		Command:     command,
		Tags:        tags,
		TimeoutSecs: timeoutSecs,
		RunAt:       runAt,
		Output:      output,
		Error:       err,
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
		ElapsedSecs: elapsedSecs,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
