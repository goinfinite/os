package scheduledTaskInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestScheduledTaskQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.ReadPersistentDbSvc()
	scheduledTaskCmdRepo := NewScheduledTaskCmdRepo(persistentDbSvc)
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	// Note: Setup/teardown are intentionally inline — test independence
	// requires each file to own its preconditions, even if it duplicates code.
	t.Run("ReadScheduledTasks", func(t *testing.T) {
		name, _ := valueObject.NewScheduledTaskName("queryTest")
		command, _ := tkValueObject.NewUnixCommand("echo scheduledTaskQueryTest")
		tag, _ := valueObject.NewScheduledTaskTag("query")
		tags := []valueObject.ScheduledTaskTag{tag}
		timeoutSecs := uint16(60)
		createDto := dto.NewCreateScheduledTask(
			name, command, tags, &timeoutSecs, nil,
		)

		taskId, err := scheduledTaskCmdRepo.Create(createDto)
		if err != nil {
			t.Fatalf("ScheduledTaskCreationFailed: %v", err)
		}
		t.Cleanup(func() {
			deleteErr := scheduledTaskCmdRepo.Delete(taskId)
			if deleteErr != nil {
				t.Errorf("ScheduledTaskCleanupFailed: %v", deleteErr)
			}
		})

		responseDto, err := scheduledTaskQueryRepo.Read(dto.ReadScheduledTasksRequest{
			Pagination: useCase.ScheduledTasksDefaultPagination,
		})
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %v", err)
		}
		if len(responseDto.Tasks) == 0 {
			t.Fatal("NoScheduledTasksFound")
		}

		isCreatedTaskReturned := false
		for _, task := range responseDto.Tasks {
			if task.Id == taskId {
				isCreatedTaskReturned = true
				break
			}
		}
		if !isCreatedTaskReturned {
			t.Errorf("CreatedTaskNotFoundInReadResponse: %s", taskId.String())
		}
	})
}
