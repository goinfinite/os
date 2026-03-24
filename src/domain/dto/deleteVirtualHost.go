package dto

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type DeleteVirtualHost struct {
	Hostname          tkValueObject.Fqdn      `json:"hostname"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

func NewDeleteVirtualHost(
	hostname tkValueObject.Fqdn,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteVirtualHost {
	return DeleteVirtualHost{
		Hostname:          hostname,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
