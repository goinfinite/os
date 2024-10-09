package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
)

type ActivityRecordCmdRepo interface {
	Create(createDto dto.CreateActivityRecord) error
	Delete(deleteDto dto.DeleteActivityRecord) error
}
