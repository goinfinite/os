package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func AddCron(
	cronCmdRepo repository.CronCmdRepo,
	addCron dto.AddCron,
) error {
	err := cronCmdRepo.Add(addCron)
	if err != nil {
		log.Printf("AddCronError: %s", err)
		return errors.New("AddCronInfraError")
	}

	cronCmdLimitStr := len(addCron.Command.String())
	if cronCmdLimitStr > 75 {
		cronCmdLimitStr = 75
	}
	cronCmdShortVersion := addCron.Command.String()[:cronCmdLimitStr]
	cronLine := addCron.Schedule.String() + " " + cronCmdShortVersion

	log.Printf("Cron '%v' added.", cronLine)

	return nil
}
