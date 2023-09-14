package repository

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/entity"
)

type CronCmdRepo interface {
	Add(addCron dto.AddCron) error
	Update(cronjob entity.Cron, updateCron dto.UpdateCron) error
}
