package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type VirtualHostQueryRepo interface {
	Read() ([]entity.VirtualHost, error)
	ReadByHostname(hostname valueObject.Fqdn) (entity.VirtualHost, error)
}
