package infra

import (
	"github.com/speedianet/os/src/domain/entity"
	o11yInfra "github.com/speedianet/os/src/infra/o11y"
)

type O11yQueryRepo struct {
}

func (repo O11yQueryRepo) GetOverview() (entity.O11yOverview, error) {
	getOverviewRepo := o11yInfra.GetOverview{}
	return getOverviewRepo.Get()
}
