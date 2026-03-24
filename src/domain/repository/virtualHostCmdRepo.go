package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type VirtualHostCmdRepo interface {
	Create(dto.CreateVirtualHost) error
	Update(dto.UpdateVirtualHost) error
	Delete(tkValueObject.Fqdn) error
}
