package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type MappingCmdRepo interface {
	Create(dto.CreateMapping) (valueObject.MappingId, error)
	Delete(valueObject.MappingId) error
}
