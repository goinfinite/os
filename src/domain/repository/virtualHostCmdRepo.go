package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
)

type VirtualHostCmdRepo interface {
	Add(addDto dto.AddVirtualHost) error
	Delete(vhost entity.VirtualHost) error
}
