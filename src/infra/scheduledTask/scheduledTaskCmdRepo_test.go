package scheduledTaskInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
)

func TestScheduledTaskCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	scheduledTaskCmdRepo := NewScheduledTaskCmdRepo(persistentDbSvc)
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	t.Run("CreateScheduledTask", func(t *testing.T) {
		name, _ := valueObject.NewScheduledTaskName("test")
		command, _ := valueObject.NewUnixCommand(
			infraEnvs.InfiniteOsBinary + " account get",
		)
		tag, _ := valueObject.NewScheduledTaskTag("account")
		tags := []valueObject.ScheduledTaskTag{tag}
		timeoutSecs := uint16(60)
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
		scheduledTasks, err := readScheduledTasks()
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
		readDto := dto.ReadScheduledTasksRequest{
			Pagination: useCase.ScheduledTasksDefaultPagination,
			TaskStatus: &pendingStatus,
		}

		responseDto, err := scheduledTaskQueryRepo.Read(readDto)
		if err != nil {
			t.Error(err)
			return
		}
		if len(responseDto.Tasks) == 0 {
			t.Error("NoPendingTasksFound")
			return
		}

		err = scheduledTaskCmdRepo.Run(responseDto.Tasks[0])
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}

		readDto = dto.ReadScheduledTasksRequest{
			Pagination: useCase.ScheduledTasksDefaultPagination,
			TaskId:     &responseDto.Tasks[0].Id,
		}

		responseDto, err = scheduledTaskQueryRepo.Read(readDto)
		if err != nil {
			t.Error(err)
			return
		}

		if len(responseDto.Tasks) == 0 {
			t.Error("NoTaskFound")
			return
		}

		completedTask := responseDto.Tasks[0]

		if completedTask.Status.String() != "completed" {
			t.Errorf("ExpectedCompletedButGot: %v", completedTask.Status.String())
		}

		err = scheduledTaskCmdRepo.Delete(completedTask.Id)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
