package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type RuntimeQueryRepo interface {
	ReadPhpVersionsInstalled() ([]valueObject.PhpVersion, error)
	ReadPhpConfigs(hostname tkValueObject.Fqdn) (entity.PhpConfigs, error)
}
