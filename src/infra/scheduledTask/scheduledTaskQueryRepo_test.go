package scheduledTaskInfra

import (
	"errors"
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

func getScheduledTasks() ([]entity.ScheduledTask, error) {
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	scheduledTasks, err := scheduledTaskQueryRepo.Read()
	if err != nil || len(scheduledTasks) == 0 {
		return nil, errors.New("NoScheduledTasksFound")
	}

	return scheduledTasks, nil
}

func TestScheduledTaskQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	scheduledTaskQueryRepo := NewScheduledTaskQueryRepo(persistentDbSvc)

	t.Run("ReadScheduledTasks", func(t *testing.T) {
		_, err := scheduledTaskQueryRepo.Read()
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("ReadScheduledTaskById", func(t *testing.T) {
		scheduledTasks, err := getScheduledTasks()
		if err != nil {
			t.Error(err)
			return
		}

		_, err = scheduledTaskQueryRepo.ReadById(scheduledTasks[0].Id)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("ReadScheduledTasksByStatus", func(t *testing.T) {
		pendingStatus, _ := valueObject.NewScheduledTaskStatus("pending")
		_, err := scheduledTaskQueryRepo.ReadByStatus(pendingStatus)
		if err != nil {
			t.Error(err)
			return
		}
	})
}
