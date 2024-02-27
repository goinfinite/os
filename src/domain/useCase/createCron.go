package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CreateCron(
	cronCmdRepo repository.CronCmdRepo,
	createCron dto.CreateCron,
) error {
	err := cronCmdRepo.Create(createCron)
	if err != nil {
		log.Printf("CreateCronError: %s", err)
		return errors.New("CreateCronInfraError")
	}

	cronCmdLimitStr := len(createCron.Command.String())
	if cronCmdLimitStr > 75 {
		cronCmdLimitStr = 75
	}
	cronCmdShortVersion := createCron.Command.String()[:cronCmdLimitStr]
	cronLine := createCron.Schedule.String() + " " + cronCmdShortVersion

	log.Printf("Cron '%v' created.", cronLine)

	return nil
}
