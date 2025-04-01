package dbModel

import (
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ScheduledTask struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Status      string `gorm:"not null,index"`
	Command     string `gorm:"not null"`
	Tags        []ScheduledTaskTag
	TimeoutSecs *uint16
	RunAt       *time.Time
	Output      *string
	Error       *string
	StartedAt   *time.Time
	FinishedAt  *time.Time
	ElapsedSecs *uint32
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (ScheduledTask) TableName() string {
	return "scheduled_tasks"
}

func NewScheduledTask(
	id uint64,
	name, status, command string,
	tags []ScheduledTaskTag,
	timeoutSecs *uint16,
	runAt *time.Time,
	output, err *string,
	startedAt, finishedAt *time.Time,
	elapsedSecs *uint32,
) ScheduledTask {
	model := ScheduledTask{
		Name:        name,
		Status:      status,
		Command:     command,
		TimeoutSecs: timeoutSecs,
		Tags:        tags,
		RunAt:       runAt,
		Output:      output,
		Error:       err,
		StartedAt:   startedAt,
		FinishedAt:  finishedAt,
		ElapsedSecs: elapsedSecs,
	}

	if id != 0 {
		model.ID = id
	}

	return model
}

func (model ScheduledTask) ToEntity() (taskEntity entity.ScheduledTask, err error) {
	id, err := valueObject.NewScheduledTaskId(model.ID)
	if err != nil {
		return taskEntity, err
	}

	name, err := valueObject.NewScheduledTaskName(model.Name)
	if err != nil {
		return taskEntity, err
	}

	status, err := valueObject.NewScheduledTaskStatus(model.Status)
	if err != nil {
		return taskEntity, err
	}

	command, err := valueObject.NewUnixCommand(model.Command)
	if err != nil {
		return taskEntity, err
	}

	tags := []valueObject.ScheduledTaskTag{}
	for _, rawTag := range model.Tags {
		tag, err := valueObject.NewScheduledTaskTag(rawTag.Tag)
		if err != nil {
			return taskEntity, err
		}
		tags = append(tags, tag)
	}

	var runAtPtr *valueObject.UnixTime
	if model.RunAt != nil {
		runAt := valueObject.NewUnixTimeWithGoTime(*model.RunAt)
		runAtPtr = &runAt
	}

	var outputPtr *valueObject.ScheduledTaskOutput
	if model.Output != nil {
		output, err := valueObject.NewScheduledTaskOutput(*model.Output)
		if err != nil {
			return taskEntity, err
		}
		outputPtr = &output
	}

	var taskErrorPtr *valueObject.ScheduledTaskOutput
	if model.Error != nil {
		taskError, err := valueObject.NewScheduledTaskOutput(*model.Error)
		if err != nil {
			return taskEntity, err
		}
		taskErrorPtr = &taskError
	}

	var startedAtPtr *valueObject.UnixTime
	if model.StartedAt != nil {
		startedAt := valueObject.NewUnixTimeWithGoTime(*model.StartedAt)
		startedAtPtr = &startedAt
	}

	var finishedAtPtr *valueObject.UnixTime
	if model.FinishedAt != nil {
		finishedAt := valueObject.NewUnixTimeWithGoTime(*model.FinishedAt)
		finishedAtPtr = &finishedAt
	}

	createdAt := valueObject.NewUnixTimeWithGoTime(model.CreatedAt)
	updatedAt := valueObject.NewUnixTimeWithGoTime(model.UpdatedAt)

	return entity.NewScheduledTask(
		id, name, status, command, tags, model.TimeoutSecs, runAtPtr, outputPtr,
		taskErrorPtr, startedAtPtr, finishedAtPtr, model.ElapsedSecs, createdAt, updatedAt,
	), nil
}
