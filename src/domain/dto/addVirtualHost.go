package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddVirtualHost struct {
	Hostname       valueObject.Fqdn            `json:"hostname"`
	Type           valueObject.VirtualHostType `json:"type"`
	ParentHostname *valueObject.Fqdn           `json:"parentHostname"`
}

func NewAddVirtualHost(
	hostname valueObject.Fqdn,
	vhostType valueObject.VirtualHostType,
	parentHostname *valueObject.Fqdn,
) AddVirtualHost {
	return AddVirtualHost{
		Hostname:       hostname,
		Type:           vhostType,
		ParentHostname: parentHostname,
	}
}
