package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ScheduledTaskCmdRepo interface {
	Create(dto.CreateScheduledTask) (valueObject.ScheduledTaskId, error)
	Update(updateDto dto.UpdateScheduledTask) error
	Run(pendingTask entity.ScheduledTask) error
}
