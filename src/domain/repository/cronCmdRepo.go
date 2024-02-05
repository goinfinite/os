package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type CronCmdRepo interface {
	Add(addCron dto.CreateCron) error
	Update(updateCron dto.UpdateCron) error
	Delete(cronId valueObject.CronId) error
}
