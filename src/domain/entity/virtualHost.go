package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type VirtualHost struct {
	Hostname         valueObject.Fqdn            `json:"hostname"`
	Type             valueObject.VirtualHostType `json:"type"`
	RootDirectory    valueObject.UnixFilePath    `json:"rootDirectory"`
	ParentHostname   *valueObject.Fqdn           `json:"parentHostname"`
	IsPrimary        bool                        `json:"isPrimary"`
	IsWildcard       bool                        `json:"isWildcard"`
	AliasesHostnames []valueObject.Fqdn          `json:"aliasesHostnames"`
	CreatedAt        valueObject.UnixTime        `json:"createdAt"`
}

func NewVirtualHost(
	hostname valueObject.Fqdn,
	vhostType valueObject.VirtualHostType,
	rootDirectory valueObject.UnixFilePath,
	parentHostname *valueObject.Fqdn,
	isPrimary bool,
	isWildcard bool,
	aliasesHostnames []valueObject.Fqdn,
	createdAt valueObject.UnixTime,
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
