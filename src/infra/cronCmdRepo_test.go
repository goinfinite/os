package infra

import (
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func addDummyCron() error {
	schedule := valueObject.NewCronSchedulePanic("* * * * *")
	command := valueObject.NewUnixCommandPanic("echo \"cronTest\" >> crontab_log.txt")
	comment := valueObject.NewCronCommentPanic("Test cron job")

	addCron := dto.AddCron{
		Schedule: schedule,
		Command:  command,
		Comment:  &comment,
	}

	cronCmdRepo := CronCmdRepo{}
	err := cronCmdRepo.Add(addCron)
	if err != nil {
		return err
	}

	return nil
}

func TestCronCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("AddCron", func(t *testing.T) {
		err := addDummyCron()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateCron", func(t *testing.T) {
		cronQueryRepo := CronQueryRepo{}
		cronCmdRepo := CronCmdRepo{}

		cron, err := cronQueryRepo.GetById(valueObject.NewCronIdPanic((1)))
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		schedule := valueObject.NewCronSchedulePanic("* * * * 0")
		command := valueObject.NewUnixCommandPanic("echo \"cronUpdateTest\" >> crontab_logs.txt")
		comment := valueObject.NewCronCommentPanic("update test")

		updateCron := dto.NewUpdateCron(
			cron.Id,
			&schedule,
			&command,
			&comment,
		)

		err = cronCmdRepo.Update(cron, updateCron)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("DeleteCron", func(t *testing.T) {
		cronCmdRepo := CronCmdRepo{}

		err := cronCmdRepo.Delete(valueObject.NewCronIdPanic((1)))
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
