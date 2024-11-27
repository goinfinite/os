package cronInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestCronCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	cronCmdRepo, err := NewCronCmdRepo()
	if err != nil {
		t.Errorf("UnexpectedError: %v", err)
	}

	ipAddress := valueObject.NewLocalhostIpAddress()
	operatorAccountId, _ := valueObject.NewAccountId(0)

	t.Run("CreateCron", func(t *testing.T) {
		schedule, _ := valueObject.NewCronSchedule("* * * * *")
		command, _ := valueObject.NewUnixCommand("echo \"cronTest\" >> crontab_log.txt")
		comment, _ := valueObject.NewCronComment("Test cron job")
		createCron := dto.NewCreateCron(
			schedule, command, &comment, operatorAccountId, ipAddress,
		)

		_, err = cronCmdRepo.Create(createCron)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateCron", func(t *testing.T) {
		id, _ := valueObject.NewCronId(1)
		schedule, _ := valueObject.NewCronSchedule("* * * * 0")
		command, _ := valueObject.NewUnixCommand("echo \"cronUpdateTest\" >> crontab_logs.txt")
		comment, _ := valueObject.NewCronComment("update test")
		updateCron := dto.NewUpdateCron(
			id, &schedule, &command, &comment, operatorAccountId, ipAddress,
		)

		err = cronCmdRepo.Update(updateCron)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("DeleteCron", func(t *testing.T) {
		id, _ := valueObject.NewCronId(1)
		err = cronCmdRepo.Delete(id)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
