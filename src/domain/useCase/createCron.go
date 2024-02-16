package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CreateCron(
	cronCmdRepo repository.CronCmdRepo,
	addCron dto.CreateCron,
) error {
	err := cronCmdRepo.Create(addCron)
	if err != nil {
		log.Printf("CreateCronError: %s", err)
		return errors.New("CreateCronInfraError")
	}

	cronCmdLimitStr := len(addCron.Command.String())
	if cronCmdLimitStr > 75 {
		cronCmdLimitStr = 75
	}
	cronCmdShortVersion := addCron.Command.String()[:cronCmdLimitStr]
	cronLine := addCron.Schedule.String() + " " + cronCmdShortVersion

	log.Printf("Cron '%v' created.", cronLine)

	return nil
}
