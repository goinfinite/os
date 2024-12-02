package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type CronQueryRepo interface {
	Read(dto.ReadCronsRequest) (dto.ReadCronsResponse, error)
	ReadFirst(dto.ReadCronsRequest) (entity.Cron, error)
}
