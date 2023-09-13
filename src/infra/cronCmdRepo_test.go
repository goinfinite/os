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
}
