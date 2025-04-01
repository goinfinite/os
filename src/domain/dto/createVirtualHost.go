package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateVirtualHost struct {
	Hostname          valueObject.Fqdn            `json:"hostname"`
	Type              valueObject.VirtualHostType `json:"type"`
	IsWildcard        *bool                       `json:"isWildcard"`
	ParentHostname    *valueObject.Fqdn           `json:"parentHostname"`
	OperatorAccountId valueObject.AccountId       `json:"-"`
	OperatorIpAddress valueObject.IpAddress       `json:"-"`
}

func NewCreateVirtualHost(
	hostname valueObject.Fqdn,
	vhostType valueObject.VirtualHostType,
	isWildcard *bool,
	parentHostname *valueObject.Fqdn,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateVirtualHost {
	return CreateVirtualHost{
		Hostname:          hostname,
		Type:              vhostType,
		IsWildcard:        isWildcard,
		ParentHostname:    parentHostname,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
