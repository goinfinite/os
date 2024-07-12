package dbModel

import (
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ScheduledTask struct {
	ID          uint   `gorm:"primarykey"`
	Name        string `gorm:"not null"`
	Status      string `gorm:"not null"`
	Command     string `gorm:"not null"`
	Tags        *string
	TimeoutSecs *uint
	RunAt       *time.Time
	Output      *string
	Error       *string
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (ScheduledTask) TableName() string {
	return "scheduled_tasks"
}

func NewScheduledTask(
	id uint,
	name, status, command string,
	tags []string,
	timeoutSecs *uint,
	runAt *time.Time,
	output, err *string,
) ScheduledTask {
	model := ScheduledTask{
		Name:        name,
		Status:      status,
		Command:     command,
		TimeoutSecs: timeoutSecs,
		RunAt:       runAt,
		Output:      output,
		Error:       err,
	}

	if id != 0 {
		model.ID = id
	}

	if len(tags) > 0 {
		modelTags := strings.Join(tags, ";")
		model.Tags = &modelTags
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
	if model.Tags != nil {
		tagsParts := strings.Split(*model.Tags, ";")
		for _, tagPart := range tagsParts {
			tag, err := valueObject.NewScheduledTaskTag(tagPart)
			if err != nil {
				return taskEntity, err
			}
			tags = append(tags, tag)
		}
	}

	var timeoutSecs *uint
	if model.TimeoutSecs != nil {
		timeoutSecs = model.TimeoutSecs
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

	createdAt := valueObject.NewUnixTimeWithGoTime(model.CreatedAt)
	updatedAt := valueObject.NewUnixTimeWithGoTime(model.UpdatedAt)

	return entity.NewScheduledTask(
		id,
		name,
		status,
		command,
		tags,
		timeoutSecs,
		runAtPtr,
		outputPtr,
		taskErrorPtr,
		createdAt,
		updatedAt,
	), nil
}
