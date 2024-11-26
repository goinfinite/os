package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type CronCmdRepo interface {
	Create(dto.CreateCron) (valueObject.CronId, error)
	Update(dto.UpdateCron) error
	Delete(valueObject.CronId) error
}
