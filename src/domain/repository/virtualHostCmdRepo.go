package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type VirtualHostCmdRepo interface {
	Create(createDto dto.CreateVirtualHost) error
	Delete(vhost entity.VirtualHost) error
	CreateMapping(createMapping dto.CreateMapping) error
	DeleteMapping(mapping entity.Mapping) error
	DeleteAutoMapping(serviceName valueObject.ServiceName) error
	RecreateMapping(mapping entity.Mapping) error
}
