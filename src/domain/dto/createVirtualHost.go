package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateVirtualHost struct {
	Hostname          valueObject.Fqdn            `json:"hostname"`
	Type              valueObject.VirtualHostType `json:"type"`
	ParentHostname    *valueObject.Fqdn           `json:"parentHostname"`
	OperatorAccountId valueObject.AccountId       `json:"-"`
	OperatorIpAddress valueObject.IpAddress       `json:"-"`
}

func NewCreateVirtualHost(
	hostname valueObject.Fqdn,
	vhostType valueObject.VirtualHostType,
	parentHostname *valueObject.Fqdn,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateVirtualHost {
	return CreateVirtualHost{
		Hostname:          hostname,
		Type:              vhostType,
		ParentHostname:    parentHostname,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
