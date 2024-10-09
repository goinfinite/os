package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
)

type ScheduledTaskQueryRepo interface {
	Read(dto.ReadScheduledTasksRequest) (dto.ReadScheduledTasksResponse, error)
}
