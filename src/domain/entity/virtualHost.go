package entity

import "github.com/speedianet/os/src/domain/valueObject"

type VirtualHost struct {
	Id             valueObject.VirtualHostId   `json:"id"`
	Hostname       valueObject.Fqdn            `json:"hostname"`
	Type           valueObject.VirtualHostType `json:"type"`
	RootDirectory  valueObject.UnixFilePath    `json:"rootDirectory"`
	ParentHostname *valueObject.Fqdn           `json:"parentHostname"`
	Mappings       []Mapping                   `json:"mappings"`
}

func NewVirtualHost(
	id valueObject.VirtualHostId,
	hostname valueObject.Fqdn,
	vhostType valueObject.VirtualHostType,
	rootDirectory valueObject.UnixFilePath,
	parentHostname *valueObject.Fqdn,
	mappings []Mapping,
) VirtualHost {
	return VirtualHost{
		Id:             id,
		Hostname:       hostname,
		Type:           vhostType,
		RootDirectory:  rootDirectory,
		ParentHostname: parentHostname,
		Mappings:       mappings,
	}
}
