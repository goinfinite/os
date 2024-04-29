package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MappingQueryRepo interface {
	GetById(id valueObject.MappingId) (entity.Mapping, error)
	GetByHostname(hostname valueObject.Fqdn) ([]entity.Mapping, error)
	GetByServiceName(serviceName valueObject.ServiceName) ([]entity.Mapping, error)
}
