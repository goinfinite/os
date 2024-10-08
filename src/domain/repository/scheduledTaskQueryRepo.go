package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ScheduledTaskQueryRepo interface {
	Read() ([]entity.ScheduledTask, error)
	ReadById(id valueObject.ScheduledTaskId) (entity.ScheduledTask, error)
	ReadByStatus(status valueObject.ScheduledTaskStatus) ([]entity.ScheduledTask, error)
}
