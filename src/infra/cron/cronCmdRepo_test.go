package cronInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestCronCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("AddCron", func(t *testing.T) {
		schedule := valueObject.NewCronSchedulePanic("* * * * *")
		command := valueObject.NewUnixCommandPanic("echo \"cronTest\" >> crontab_log.txt")
		comment := valueObject.NewCronCommentPanic("Test cron job")

		addCron := dto.NewAddCron(
			schedule,
			command,
			&comment,
		)

		cronCmdRepo, err := NewCronCmdRepo()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		err = cronCmdRepo.Add(addCron)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateCron", func(t *testing.T) {
		cronCmdRepo, err := NewCronCmdRepo()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		schedule := valueObject.NewCronSchedulePanic("* * * * 0")
		command := valueObject.NewUnixCommandPanic("echo \"cronUpdateTest\" >> crontab_logs.txt")
		comment := valueObject.NewCronCommentPanic("update test")

		updateCron := dto.NewUpdateCron(
			valueObject.NewCronIdPanic(1),
			&schedule,
			&command,
			&comment,
		)

		err = cronCmdRepo.Update(updateCron)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("DeleteCron", func(t *testing.T) {
		cronCmdRepo, err := NewCronCmdRepo()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		err = cronCmdRepo.Delete(valueObject.NewCronIdPanic((1)))
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
