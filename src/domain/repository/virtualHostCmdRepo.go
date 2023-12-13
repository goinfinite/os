package repository

import (
	"github.com/speedianet/os/src/domain/dto"
)

type VirtualHostCmdRepo interface {
	Add(addDto dto.AddVirtualHost) error
}
