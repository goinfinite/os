package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type ScheduledTask struct {
	Id          valueObject.ScheduledTaskId      `json:"id"`
	Name        valueObject.ScheduledTaskName    `json:"name"`
	Status      valueObject.ScheduledTaskStatus  `json:"status"`
	Command     valueObject.UnixCommand          `json:"command"`
	Tags        []valueObject.ScheduledTaskTag   `json:"tags"`
	TimeoutSecs *uint16                          `json:"timeoutSecs"`
	RunAt       *valueObject.UnixTime            `json:"runAt"`
	Output      *valueObject.ScheduledTaskOutput `json:"output"`
	Error       *valueObject.ScheduledTaskOutput `json:"err"`
	StartedAt   *valueObject.UnixTime            `json:"startedAt"`
	FinishedAt  *valueObject.UnixTime            `json:"finishedAt"`
	ElapsedSecs *uint32                          `json:"elapsedSecs"`
	CreatedAt   valueObject.UnixTime             `json:"createdAt"`
	UpdatedAt   valueObject.UnixTime             `json:"updatedAt"`
}

func NewScheduledTask(
	id valueObject.ScheduledTaskId,
	name valueObject.ScheduledTaskName,
	status valueObject.ScheduledTaskStatus,
	command valueObject.UnixCommand,
	tags []valueObject.ScheduledTaskTag,
	timeoutSecs *uint16,
	runAt *valueObject.UnixTime,
	output, err *valueObject.ScheduledTaskOutput,
	startedAt, finishedAt *valueObject.UnixTime,
	elapsedSecs *uint32,
	createdAt, updatedAt valueObject.UnixTime,
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
