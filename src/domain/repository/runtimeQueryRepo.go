package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type RuntimeQueryRepo interface {
	GetPhpVersionsInstalled() ([]valueObject.PhpVersion, error)
	GetPhpConfigs(hostname valueObject.Fqdn) (entity.PhpConfigs, error)
}
