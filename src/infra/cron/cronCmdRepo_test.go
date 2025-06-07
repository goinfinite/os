package cronInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestCronCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	cronCmdRepo := NewCronCmdRepo()

	var id valueObject.CronId
	schedule, _ := valueObject.NewCronSchedule("* * * * *")
	command, _ := valueObject.NewUnixCommand("echo \"cronTest\" >> crontab_log.txt")
	comment, _ := valueObject.NewCronComment("Test cron job")
	operatorAccountId := valueObject.AccountIdSystem
	operatorIpAddress := valueObject.IpAddressSystem

	createCron := dto.NewCreateCron(
		schedule, command, &comment, operatorAccountId, operatorIpAddress,
	)

	t.Run("CreateCron", func(t *testing.T) {
		var err error
		id, err = cronCmdRepo.Create(createCron)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: '%s'", err.Error())
		}
	})

	t.Run("UpdateCron", func(t *testing.T) {
		schedule, _ = valueObject.NewCronSchedule("* * * * 0")
		updateCron := dto.NewUpdateCron(
			id, &schedule, nil, nil, []string{}, operatorAccountId, operatorIpAddress,
		)

		err := cronCmdRepo.Update(updateCron)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: '%s'", err.Error())
		}
	})

	t.Run("DeleteCron", func(t *testing.T) {
		err := cronCmdRepo.Delete(id)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: '%s'", err.Error())
		}
	})
}
