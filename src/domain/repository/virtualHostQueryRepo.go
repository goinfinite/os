package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type VirtualHostQueryRepo interface {
	GetVirtualHostConfFilePath(vhost valueObject.Fqdn) (valueObject.UnixFilePath, error)
	Get() ([]entity.VirtualHost, error)
	GetByHostname(hostname valueObject.Fqdn) (entity.VirtualHost, error)
	GetWithMappings() ([]dto.VirtualHostWithMappings, error)
	GetMappingById(
		vhostHostname valueObject.Fqdn,
		id valueObject.MappingId,
	) (entity.Mapping, error)
}
