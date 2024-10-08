package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func UpdateCron(
	cronQueryRepo repository.CronQueryRepo,
	cronCmdRepo repository.CronCmdRepo,
	updateCron dto.UpdateCron,
) error {
	_, err := cronQueryRepo.ReadById(updateCron.Id)
	if err != nil {
		log.Printf("CronNotFound: %s", err)
		return errors.New("CronNotFound")
	}

	err = cronCmdRepo.Update(updateCron)
	if err != nil {
		log.Printf("UpdateCronError: %s", err)
		return errors.New("UpdateCronInfraError")
	}

	log.Printf("Cron with ID '%v' updated.", updateCron.Id.String())

	return nil
}
