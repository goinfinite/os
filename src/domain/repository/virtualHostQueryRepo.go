package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type VirtualHostQueryRepo interface {
	Get() ([]entity.VirtualHost, error)
	GetByHostname(hostname valueObject.Fqdn) (entity.VirtualHost, error)
	GetWithMappings() ([]dto.VirtualHostWithMappings, error)
	GetMappingsByHostname(
		hostname valueObject.Fqdn,
	) ([]entity.Mapping, error)
	GetMappingById(
		vhostHostname valueObject.Fqdn,
		id valueObject.MappingId,
	) (entity.Mapping, error)
	IsDomainOwner(vhost valueObject.Fqdn, ownershipHash string) bool
}
