package dto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type VirtualHostWithMappings struct {
	entity.VirtualHost
	Mappings []valueObject.Mapping `json:"mappings,omitempty"`
}

func NewVirtualHostWithMappings(
	vhost entity.VirtualHost,
	mappings []valueObject.Mapping,
) VirtualHostWithMappings {
	return VirtualHostWithMappings{
		VirtualHost: vhost,
		Mappings:    mappings,
	}
}
