package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ReadVirtualHostsRequest struct {
	Pagination       tkDto.Pagination                     `json:"pagination"`
	Hostname         *tkValueObject.Fqdn                  `json:"hostname"`
	VirtualHostType  *valueObject.VirtualHostType         `json:"type"`
	RootDirectory    *tkValueObject.UnixAbsoluteFilePath  `json:"rootDirectory"`
	ParentHostname   *tkValueObject.Fqdn                  `json:"parentHostname"`
	WithMappings     *bool                                `json:"withMappings"`
	IsWildcard       *bool                                `json:"isWildcard"`
	IsPrimary        *bool                                `json:"-"`
	AliasesHostnames []tkValueObject.Fqdn                 `json:"aliasesHostnames"`
	CreatedBeforeAt  *tkValueObject.UnixTime              `json:"createdBeforeAt"`
	CreatedAfterAt   *tkValueObject.UnixTime              `json:"createdAfterAt"`
}

type VirtualHostWithMappings struct {
	entity.VirtualHost
	Mappings []entity.Mapping `json:"mappings"`
}

type ReadVirtualHostsResponse struct {
	Pagination              tkDto.Pagination          `json:"pagination"`
	VirtualHosts            []entity.VirtualHost      `json:"virtualHosts"`
	VirtualHostWithMappings []VirtualHostWithMappings `json:"virtualHostWithMappings"`
}
