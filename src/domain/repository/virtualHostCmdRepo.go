package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type VirtualHostCmdRepo interface {
	Create(dto.CreateVirtualHost) error
	Delete(valueObject.Fqdn) error
}
