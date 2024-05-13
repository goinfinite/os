package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type VirtualHostQueryRepo interface {
	Get() ([]entity.VirtualHost, error)
	GetByHostname(hostname valueObject.Fqdn) (entity.VirtualHost, error)
}
