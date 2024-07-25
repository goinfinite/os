package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type CronQueryRepo interface {
	Read() ([]entity.Cron, error)
	ReadById(cronId valueObject.CronId) (entity.Cron, error)
}
