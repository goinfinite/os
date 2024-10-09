package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateVirtualHost struct {
	Hostname       valueObject.Fqdn            `json:"hostname"`
	Type           valueObject.VirtualHostType `json:"type"`
	ParentHostname *valueObject.Fqdn           `json:"parentHostname"`
}

func NewCreateVirtualHost(
	hostname valueObject.Fqdn,
	vhostType valueObject.VirtualHostType,
	parentHostname *valueObject.Fqdn,
) CreateVirtualHost {
	return CreateVirtualHost{
		Hostname:       hostname,
		Type:           vhostType,
		ParentHostname: parentHostname,
	}
}
