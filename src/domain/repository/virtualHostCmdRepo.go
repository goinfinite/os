package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
)

type VirtualHostCmdRepo interface {
	Add(addDto dto.CreateVirtualHost) error
	Delete(vhost entity.VirtualHost) error
	CreateMapping(addMapping dto.CreateMapping) error
	DeleteMapping(mapping entity.Mapping) error
	RecreateMapping(mapping entity.Mapping) error
}
