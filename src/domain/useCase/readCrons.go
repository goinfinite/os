package useCase

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadCrons(
	cronQueryRepo repository.CronQueryRepo,
) ([]entity.Cron, error) {
	return cronQueryRepo.Read()
}
