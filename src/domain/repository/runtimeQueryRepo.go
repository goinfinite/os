package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type RuntimeQueryRepo interface {
	ReadPhpVersionsInstalled() ([]valueObject.PhpVersion, error)
	ReadPhpConfigs(hostname valueObject.Fqdn) (entity.PhpConfigs, error)
}
