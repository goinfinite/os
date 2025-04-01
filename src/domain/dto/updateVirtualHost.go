package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateVirtualHost struct {
	Hostname          valueObject.Fqdn      `json:"hostname"`
	IsWildcard        *bool                 `json:"isWildcard"`
	OperatorAccountId valueObject.AccountId `json:"-"`
	OperatorIpAddress valueObject.IpAddress `json:"-"`
}

func NewUpdateVirtualHost(
	hostname valueObject.Fqdn,
	isWildcard *bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateVirtualHost {
	return UpdateVirtualHost{
		Hostname:          hostname,
		IsWildcard:        isWildcard,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
