package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type VirtualHostCmdRepo interface {
	Create(createDto dto.CreateVirtualHost) error
	Delete(vhost entity.VirtualHost) error
}
