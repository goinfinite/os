package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type VirtualHostQueryRepo interface {
	Read(dto.ReadVirtualHostsRequest) (dto.ReadVirtualHostsResponse, error)
	ReadFirst(dto.ReadVirtualHostsRequest) (entity.VirtualHost, error)
	ReadFirstWithMappings(dto.ReadVirtualHostsRequest) (dto.VirtualHostWithMappings, error)
}
