package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type RuntimeQueryRepo interface {
	GetPhpVersionsInstalled() ([]valueObject.PhpVersion, error)
	GetPhpConfigs(hostname valueObject.Fqdn) (entity.PhpConfigs, error)
}
