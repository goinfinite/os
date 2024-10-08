package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type ActivityRecordQueryRepo interface {
	Read(readDto dto.ReadActivityRecords) ([]entity.ActivityRecord, error)
}
