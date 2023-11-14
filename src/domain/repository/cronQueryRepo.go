package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type CronQueryRepo interface {
	Get() ([]entity.Cron, error)
	GetById(cronId valueObject.CronId) (entity.Cron, error)
}
