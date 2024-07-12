package scheduledTaskInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestScheduledTaskCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	scheduledTaskCmdRepo := NewScheduledTaskCmdRepo(persistentDbSvc)
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	t.Run("CreateScheduledTask", func(t *testing.T) {
		name, _ := valueObject.NewScheduledTaskName("test")
		command, _ := valueObject.NewUnixCommand("/var/speedia/control account get")
		containerTag, _ := valueObject.NewScheduledTaskTag("account")
		tags := []valueObject.ScheduledTaskTag{containerTag}
		timeoutSecs := uint(60)
		runAt := valueObject.NewUnixTimeNow()

		createDto := dto.NewCreateScheduledTask(
			name, command, tags, &timeoutSecs, &runAt,
		)

		err := scheduledTaskCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("UpdateScheduledTask", func(t *testing.T) {
		scheduledTasks, err := getScheduledTasks()
		if err != nil {
			t.Error(err)
			return
		}

		newStatus, _ := valueObject.NewScheduledTaskStatus("pending")
		updateDto := dto.NewUpdateScheduledTask(scheduledTasks[0].Id, &newStatus, nil)

		err = scheduledTaskCmdRepo.Update(updateDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("RunScheduledTasks", func(t *testing.T) {
		pendingStatus, _ := valueObject.NewScheduledTaskStatus("pending")
		pendingTasks, err := scheduledTaskQueryRepo.ReadByStatus(pendingStatus)
		if err != nil {
			t.Error(err)
			return
		}

		err = scheduledTaskCmdRepo.Run(pendingTasks[0])
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}

		completedTask, err := scheduledTaskQueryRepo.ReadById(pendingTasks[0].Id)
		if err != nil {
			t.Error(err)
			return
		}

		if completedTask.Status.String() != "completed" {
			t.Errorf("ExpectedCompletedButGot: %v", completedTask.Status.String())
		}

		err = scheduledTaskCmdRepo.Delete(completedTask.Id)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
