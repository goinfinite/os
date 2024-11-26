package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type CronQueryRepo interface {
	Read() ([]entity.Cron, error)
	ReadById(valueObject.CronId) (entity.Cron, error)
	ReadByComment(valueObject.CronComment) (entity.Cron, error)
}
