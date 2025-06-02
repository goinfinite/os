package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type RuntimeCmdRepo interface {
	RunPhpCommand(dto.RunPhpCommandRequest) (dto.RunPhpCommandResponse, error)
	UpdatePhpVersion(valueObject.Fqdn, valueObject.PhpVersion) error
	UpdatePhpSettings(valueObject.Fqdn, []entity.PhpSetting) error
	UpdatePhpModules(valueObject.Fqdn, []entity.PhpModule) error
}
