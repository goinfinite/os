package o11yInfra

import (
	"github.com/speedianet/os/src/domain/entity"
)

type O11yQueryRepo struct {
}

func (repo O11yQueryRepo) GetOverview() (entity.O11yOverview, error) {
	getOverviewRepo := GetOverview{}
	return getOverviewRepo.Get()
}
