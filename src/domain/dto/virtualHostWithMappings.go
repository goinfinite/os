package dto

import (
	"encoding/json"

	"github.com/speedianet/os/src/domain/entity"
)

type VirtualHostWithMappings struct {
	entity.VirtualHost
	Mappings []entity.Mapping `json:"mappings"`
}

func NewVirtualHostWithMappings(
	vhost entity.VirtualHost,
	mappings []entity.Mapping,
) VirtualHostWithMappings {
	return VirtualHostWithMappings{
		VirtualHost: vhost,
		Mappings:    mappings,
	}
}

func (dto VirtualHostWithMappings) JsonSerialize() string {
	jsonBytes, _ := json.Marshal(dto)
	return string(jsonBytes)
}
