package repository

import (
	"github.com/speedianet/os/src/domain/dto"
)

type ActivityRecordCmdRepo interface {
	Create(createDto dto.CreateActivityRecord) error
}
