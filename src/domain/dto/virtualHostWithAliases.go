package dto

import (
	"github.com/speedianet/os/src/domain/entity"
)

type VirtualHostWithAliases struct {
	entity.VirtualHost
	Aliases []entity.VirtualHost
}

func NewVirtualHostWithAliases(
	vhost entity.VirtualHost,
	aliases []entity.VirtualHost,
) VirtualHostWithAliases {
	return VirtualHostWithAliases{
		VirtualHost: vhost,
		Aliases:     aliases,
	}
}
