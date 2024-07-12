package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
)

type ScheduledTaskCmdRepo interface {
	Create(createDto dto.CreateScheduledTask) error
	Update(updateDto dto.UpdateScheduledTask) error
	Run(pendingTask entity.ScheduledTask) error
}
