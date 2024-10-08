package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type ScheduledTaskCmdRepo interface {
	Create(createDto dto.CreateScheduledTask) error
	Update(updateDto dto.UpdateScheduledTask) error
	Run(pendingTask entity.ScheduledTask) error
}
