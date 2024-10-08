package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type MappingQueryRepo interface {
	ReadById(id valueObject.MappingId) (entity.Mapping, error)
	ReadByHostname(hostname valueObject.Fqdn) ([]entity.Mapping, error)
	ReadByServiceName(serviceName valueObject.ServiceName) ([]entity.Mapping, error)
	ReadWithMappings() ([]dto.VirtualHostWithMappings, error)
}
