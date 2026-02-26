package dto

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type UpdateVirtualHost struct {
	Hostname          tkValueObject.Fqdn      `json:"hostname"`
	IsWildcard        *bool                   `json:"isWildcard"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

func NewUpdateVirtualHost(
	hostname tkValueObject.Fqdn,
	isWildcard *bool,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) UpdateVirtualHost {
	return UpdateVirtualHost{
		Hostname:          hostname,
		IsWildcard:        isWildcard,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
