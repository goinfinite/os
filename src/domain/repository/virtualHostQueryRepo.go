package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
)

type VirtualHostQueryRepo interface {
	Get() ([]entity.VirtualHost, error)
	GetWithMappings() ([]dto.VirtualHostWithMappings, error)
}
