package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MappingCmdRepo interface {
	Create(createDto dto.CreateMapping) (valueObject.MappingId, error)
	DeleteMapping(mappingId valueObject.MappingId) error
	DeleteAutoMapping(serviceName valueObject.ServiceName) error
}
