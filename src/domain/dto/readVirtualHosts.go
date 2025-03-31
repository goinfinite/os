package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadVirtualHostsRequest struct {
	Pagination       Pagination                   `json:"pagination"`
	Hostname         *valueObject.Fqdn            `json:"hostname"`
	VirtualHostType  *valueObject.VirtualHostType `json:"type"`
	RootDirectory    *valueObject.UnixFilePath    `json:"rootDirectory"`
	ParentHostname   *valueObject.Fqdn            `json:"parentHostname"`
	WithMappings     *bool                        `json:"withMappings"`
	IsWildcard       *bool                        `json:"isWildcard"`
	IsPrimary        *bool                        `json:"-"`
	AliasesHostnames []valueObject.Fqdn           `json:"aliasesHostnames"`
	CreatedBeforeAt  *valueObject.UnixTime        `json:"createdBeforeAt"`
	CreatedAfterAt   *valueObject.UnixTime        `json:"createdAfterAt"`
}

type VirtualHostWithMappings struct {
	entity.VirtualHost
	Mappings []entity.Mapping `json:"mappings"`
}

type ReadVirtualHostsResponse struct {
	Pagination              Pagination                `json:"pagination"`
	VirtualHosts            []entity.VirtualHost      `json:"virtualHosts"`
	VirtualHostWithMappings []VirtualHostWithMappings `json:"virtualHostWithMappings"`
}
