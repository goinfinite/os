package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
)

type CronQueryRepo interface {
	Get() ([]entity.Cron, error)
}
