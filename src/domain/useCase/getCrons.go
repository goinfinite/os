package useCase

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/repository"
)

func GetCrons(
	cronQueryRepo repository.CronQueryRepo,
) ([]entity.Cron, error) {
	return cronQueryRepo.Get()
}
