package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MappingQueryRepo interface {
	GetByHostname(hostname valueObject.Fqdn) ([]entity.Mapping, error)
}
