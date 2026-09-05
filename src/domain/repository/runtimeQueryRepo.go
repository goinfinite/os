package repository

import (
	"errors"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

var (
	ErrPhpModuleNotSupportedForVersion = errors.New(
		"PhpModuleNotSupportedForVersion",
	)
	ErrPhpVersionChanged = errors.New("PhpVersionChanged")
)

type RuntimeQueryRepo interface {
	ReadPhpModules(valueObject.PhpVersion) ([]entity.PhpModule, error)
	ReadPhpVersion(tkValueObject.Fqdn) (entity.PhpVersion, error)
	ReadPhpVersionsInstalled() ([]valueObject.PhpVersion, error)
	ReadPhpConfigs(hostname tkValueObject.Fqdn) (entity.PhpConfigs, error)
}
