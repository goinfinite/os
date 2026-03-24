package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type VirtualHost struct {
	Hostname         tkValueObject.Fqdn                 `json:"hostname"`
	Type             valueObject.VirtualHostType        `json:"type"`
	RootDirectory    tkValueObject.UnixAbsoluteFilePath `json:"rootDirectory"`
	ParentHostname   *tkValueObject.Fqdn                `json:"parentHostname"`
	IsPrimary        bool                               `json:"isPrimary"`
	IsWildcard       bool                               `json:"isWildcard"`
	AliasesHostnames []tkValueObject.Fqdn               `json:"aliasesHostnames"`
	CreatedAt        tkValueObject.UnixTime             `json:"createdAt"`
}

func NewVirtualHost(
	hostname tkValueObject.Fqdn,
	vhostType valueObject.VirtualHostType,
	rootDirectory tkValueObject.UnixAbsoluteFilePath,
	parentHostname *tkValueObject.Fqdn,
	isPrimary bool,
	isWildcard bool,
	aliasesHostnames []tkValueObject.Fqdn,
	createdAt tkValueObject.UnixTime,
) VirtualHost {
	return VirtualHost{
		Hostname:         hostname,
		Type:             vhostType,
		RootDirectory:    rootDirectory,
		ParentHostname:   parentHostname,
		IsPrimary:        isPrimary,
		IsWildcard:       isWildcard,
		AliasesHostnames: aliasesHostnames,
		CreatedAt:        createdAt,
	}
}
