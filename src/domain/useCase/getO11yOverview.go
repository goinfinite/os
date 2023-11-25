package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetO11yOverview(
	o11yQueryRepo repository.O11yQueryRepo,
) (entity.O11yOverview, error) {
	return o11yQueryRepo.GetOverview()
}
