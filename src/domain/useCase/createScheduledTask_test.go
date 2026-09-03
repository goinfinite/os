package useCase

import (
	"testing"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type createScheduledTaskCmdRepo struct {
	taskId valueObject.ScheduledTaskId
}

func (repo createScheduledTaskCmdRepo) Create(
	createDto dto.CreateScheduledTask,
) (valueObject.ScheduledTaskId, error) {
	return repo.taskId, nil
}

func (repo createScheduledTaskCmdRepo) Update(
	updateDto dto.UpdateScheduledTask,
) error {
	return nil
}

func (repo createScheduledTaskCmdRepo) Run(
	pendingTask entity.ScheduledTask,
) error {
	return nil
}

func TestCreateScheduledTaskReturnsTaskId(test *testing.T) {
	expectedTaskId, err := valueObject.NewScheduledTaskId(42)
	if err != nil {
		test.Fatalf("TaskIdCreationFailed: %v", err)
	}

	name, err := valueObject.NewScheduledTaskName("test")
	if err != nil {
		test.Fatalf("TaskNameCreationFailed: %v", err)
	}
	command, err := tkValueObject.NewUnixCommand("true")
	if err != nil {
		test.Fatalf("CommandCreationFailed: %v", err)
	}

	actualTaskId, err := CreateScheduledTask(
		createScheduledTaskCmdRepo{taskId: expectedTaskId},
		dto.NewCreateScheduledTask(name, command, nil, nil, nil),
	)
	if err != nil {
		test.Fatalf("CreateScheduledTaskFailed: %v", err)
	}
	if actualTaskId != expectedTaskId {
		test.Errorf(
			"TaskIdMismatch: expected %d, got %d",
			expectedTaskId.Uint64(), actualTaskId.Uint64(),
		)
	}
}
