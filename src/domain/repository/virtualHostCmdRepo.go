package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type VirtualHostCmdRepo interface {
	Add(addDto dto.AddVirtualHost) error
	Delete(vhost entity.VirtualHost) error
	AddMapping(addMapping dto.AddMapping) error
	DeleteMapping(hostname valueObject.Fqdn, mapping valueObject.Mapping) error
}
