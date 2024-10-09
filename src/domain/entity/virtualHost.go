package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type VirtualHost struct {
	Hostname       valueObject.Fqdn            `json:"hostname"`
	Type           valueObject.VirtualHostType `json:"type"`
	RootDirectory  valueObject.UnixFilePath    `json:"rootDirectory"`
	ParentHostname *valueObject.Fqdn           `json:"parentHostname"`
}

func NewVirtualHost(
	hostname valueObject.Fqdn,
	vhostType valueObject.VirtualHostType,
	rootDirectory valueObject.UnixFilePath,
	parentHostname *valueObject.Fqdn,
) VirtualHost {
	return VirtualHost{
		Hostname:       hostname,
		Type:           vhostType,
		RootDirectory:  rootDirectory,
		ParentHostname: parentHostname,
	}
}
