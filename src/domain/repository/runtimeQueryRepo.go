package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type RuntimeQueryRepo interface {
	ReadPhpVersionsInstalled() ([]valueObject.PhpVersion, error)
	ReadPhpConfigs(hostname valueObject.Fqdn) (entity.PhpConfigs, error)
	IsHtaccessModified() bool
}
