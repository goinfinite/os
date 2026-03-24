package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type RuntimeCmdRepo interface {
	RunPhpCommand(dto.RunPhpCommandRequest) (dto.RunPhpCommandResponse, error)
	UpdatePhpVersion(tkValueObject.Fqdn, valueObject.PhpVersion) error
	UpdatePhpSettings(tkValueObject.Fqdn, []entity.PhpSetting) error
	UpdatePhpModules(tkValueObject.Fqdn, []entity.PhpModule) error
}
