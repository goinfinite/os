package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type RuntimeCmdRepo interface {
	UpdatePhpVersion(hostname valueObject.Fqdn, version valueObject.PhpVersion) error
	UpdatePhpSettings(hostname valueObject.Fqdn, settings []entity.PhpSetting) error
	UpdatePhpModules(hostname valueObject.Fqdn, modules []entity.PhpModule) error
}
