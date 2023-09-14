package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type CronQueryRepo interface {
	Get() ([]entity.Cron, error)
	GetById(
		croId valueObject.CronId,
	) (entity.Cron, error)
}
