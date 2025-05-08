package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type DatabaseQueryRepo interface {
	Read(dto.ReadDatabasesRequest) (dto.ReadDatabasesResponse, error)
	ReadFirst(dto.ReadDatabasesRequest) (entity.Database, error)
}
