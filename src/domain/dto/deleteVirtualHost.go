package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteVirtualHost struct {
	Hostname           valueObject.Fqdn      `json:"hostname"`
	PrimaryVirtualHost valueObject.Fqdn      `json:"-"`
	OperatorAccountId  valueObject.AccountId `json:"-"`
	OperatorIpAddress  valueObject.IpAddress `json:"-"`
}

func NewDeleteVirtualHost(
	hostname valueObject.Fqdn,
	primaryVirtualHost valueObject.Fqdn,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteVirtualHost {
	return DeleteVirtualHost{
		Hostname:           hostname,
		PrimaryVirtualHost: primaryVirtualHost,
		OperatorAccountId:  operatorAccountId,
		OperatorIpAddress:  operatorIpAddress,
	}
}
