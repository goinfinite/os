package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type VirtualHostQueryRepo interface {
	Read() ([]entity.VirtualHost, error)
	ReadByHostname(hostname valueObject.Fqdn) (entity.VirtualHost, error)
}
