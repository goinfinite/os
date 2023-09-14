package repository

import (
	"github.com/speedianet/sam/src/domain/dto"
)

type CronCmdRepo interface {
	Add(addCron dto.AddCron) error
	Update(updateCron dto.UpdateCron) error
}
