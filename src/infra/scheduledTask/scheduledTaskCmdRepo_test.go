package scheduledTaskInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestScheduledTaskCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.ReadPersistentDbSvc()
	scheduledTaskCmdRepo := NewScheduledTaskCmdRepo(persistentDbSvc)
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	createdTaskId := valueObject.ScheduledTaskId(0)

	// Note: Setup/teardown are intentionally inline — test independence
	// requires each file to own its preconditions, even if it duplicates code.
	t.Run("CreateScheduledTask", func(t *testing.T) {
		name, _ := valueObject.NewScheduledTaskName("test")
		command, _ := tkValueObject.NewUnixCommand("echo scheduledTaskTest")
		tag, _ := valueObject.NewScheduledTaskTag("account")
		tags := []valueObject.ScheduledTaskTag{tag}
		timeoutSecs := uint16(60)
		runAt := tkValueObject.NewUnixTimeNow()

		createDto := dto.NewCreateScheduledTask(
			name, command, tags, &timeoutSecs, &runAt,
		)

		taskId, err := scheduledTaskCmdRepo.Create(createDto)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %v", err)
		}
		if taskId.Uint64() == 0 {
			t.Fatal("ExpectedGeneratedTaskId")
		}

		createdTaskId = taskId
	})
	if createdTaskId.Uint64() == 0 {
		t.Fatal("ScheduledTaskCreationFailed: dependent tests skipped")
	}
	t.Cleanup(func() {
		deleteErr := scheduledTaskCmdRepo.Delete(createdTaskId)
		if deleteErr != nil {
			t.Errorf("ScheduledTaskCleanupFailed: %v", deleteErr)
		}
	})

	t.Run("UpdateScheduledTask", func(t *testing.T) {
		completedStatus, _ := valueObject.NewScheduledTaskStatus("completed")
		updateDto := dto.NewUpdateScheduledTask(
			createdTaskId, &completedStatus, nil,
		)

		err := scheduledTaskCmdRepo.Update(updateDto)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %v", err)
		}

		readDto := dto.ReadScheduledTasksRequest{
			Pagination: useCase.ScheduledTasksDefaultPagination,
			TaskId:     &createdTaskId,
		}
		responseDto, err := scheduledTaskQueryRepo.Read(readDto)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %v", err)
		}
		if len(responseDto.Tasks) == 0 {
			t.Fatal("NoTaskFound")
		}
		if responseDto.Tasks[0].Status.String() != "completed" {
			t.Errorf(
				"ExpectedCompletedButGot: %v",
				responseDto.Tasks[0].Status.String(),
			)
		}

		pendingStatus, _ := valueObject.NewScheduledTaskStatus("pending")
		updateDto = dto.NewUpdateScheduledTask(createdTaskId, &pendingStatus, nil)

		err = scheduledTaskCmdRepo.Update(updateDto)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("RunScheduledTasks", func(t *testing.T) {
		readDto := dto.ReadScheduledTasksRequest{
			Pagination: useCase.ScheduledTasksDefaultPagination,
			TaskId:     &createdTaskId,
		}

		responseDto, err := scheduledTaskQueryRepo.Read(readDto)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %v", err)
		}
		if len(responseDto.Tasks) == 0 {
			t.Fatal("NoTaskFound")
		}

		pendingTask := responseDto.Tasks[0]
		if pendingTask.Status.String() != "pending" {
			t.Fatalf("ExpectedPendingButGot: %v", pendingTask.Status.String())
		}

		err = scheduledTaskCmdRepo.Run(pendingTask)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %v", err)
		}

		responseDto, err = scheduledTaskQueryRepo.Read(readDto)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %v", err)
		}
		if len(responseDto.Tasks) == 0 {
			t.Fatal("NoTaskFound")
		}

		completedTask := responseDto.Tasks[0]
		if completedTask.Status.String() != "completed" {
			t.Errorf("ExpectedCompletedButGot: %v", completedTask.Status.String())
		}

		err = scheduledTaskCmdRepo.Delete(completedTask.Id)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
