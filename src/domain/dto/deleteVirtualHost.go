package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteVirtualHost struct {
	PrimaryVirtualHost valueObject.Fqdn      `json:"primaryVirtualHost"`
	Hostname           valueObject.Fqdn      `json:"hostname"`
	OperatorAccountId  valueObject.AccountId `json:"-"`
	OperatorIpAddress  valueObject.IpAddress `json:"-"`
}

func NewDeleteVirtualHost(
	primaryVirtualHost valueObject.Fqdn,
	hostname valueObject.Fqdn,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteVirtualHost {
	return DeleteVirtualHost{
		PrimaryVirtualHost: primaryVirtualHost,
		Hostname:           hostname,
		OperatorAccountId:  operatorAccountId,
		OperatorIpAddress:  operatorIpAddress,
	}
}
