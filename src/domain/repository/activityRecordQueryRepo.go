package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
)

type ActivityRecordQueryRepo interface {
	Read(readDto dto.ReadActivityRecords) ([]entity.ActivityRecord, error)
}
