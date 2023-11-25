package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetCrons(
	cronQueryRepo repository.CronQueryRepo,
) ([]entity.Cron, error) {
	return cronQueryRepo.Get()
}
