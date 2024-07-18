package cronInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestCronCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	cronCmdRepo, err := NewCronCmdRepo()
	if err != nil {
		t.Errorf("UnexpectedError: %v", err)
	}

	t.Run("CreateCron", func(t *testing.T) {
		schedule, err := valueObject.NewCronSchedule("* * * * *")
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		command, err := valueObject.NewUnixCommand("echo \"cronTest\" >> crontab_log.txt")
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		comment, err := valueObject.NewCronComment("Test cron job")
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		createCron := dto.NewCreateCron(schedule, command, &comment)

		err = cronCmdRepo.Create(createCron)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateCron", func(t *testing.T) {
		id, err := valueObject.NewCronId(1)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		schedule, err := valueObject.NewCronSchedule("* * * * 0")
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		command, err := valueObject.NewUnixCommand("echo \"cronUpdateTest\" >> crontab_logs.txt")
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		comment, err := valueObject.NewCronComment("update test")
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		updateCron := dto.NewUpdateCron(id, &schedule, &command, &comment)

		err = cronCmdRepo.Update(updateCron)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("DeleteCron", func(t *testing.T) {
		id, err := valueObject.NewCronId(1)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		err = cronCmdRepo.Delete(id)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
