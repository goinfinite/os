package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type CronCmdRepo interface {
	Create(createCron dto.CreateCron) error
	Update(updateCron dto.UpdateCron) error
	Delete(cronId valueObject.CronId) error
	DeleteByComment(comment valueObject.CronComment) error
}
