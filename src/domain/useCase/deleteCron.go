package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteCron(
	cronQueryRepo repository.CronQueryRepo,
	cronCmdRepo repository.CronCmdRepo,
	cronId valueObject.CronId,
) error {
	_, err := cronQueryRepo.GetById(cronId)
	if err != nil {
		log.Printf("CronNotFound: %s", err)
		return errors.New("CronNotFound")
	}

	err = cronCmdRepo.Delete(cronId)
	if err != nil {
		log.Printf("DeleteCronError: %s", err)
		return errors.New("DeleteCronInfraError")
	}

	log.Printf("CronId '%v' deleted.", cronId)

	return nil
}
