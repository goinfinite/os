package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateVirtualHost struct {
	Hostname          tkValueObject.Fqdn          `json:"hostname"`
	Type              valueObject.VirtualHostType `json:"type"`
	IsWildcard        *bool                       `json:"isWildcard"`
	ParentHostname    *tkValueObject.Fqdn         `json:"parentHostname"`
	OperatorAccountId tkValueObject.AccountId     `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress     `json:"-"`
}

func NewCreateVirtualHost(
	hostname tkValueObject.Fqdn,
	vhostType valueObject.VirtualHostType,
	isWildcard *bool,
	parentHostname *tkValueObject.Fqdn,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
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
