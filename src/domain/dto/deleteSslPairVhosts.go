package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteSslPairVhosts struct {
	SslPairId             valueObject.SslPairId `json:"sslPairId"`
	VirtualHostsHostnames []valueObject.Fqdn    `json:"virtualHostsHostnames"`
	OperatorAccountId     valueObject.AccountId `json:"-"`
	OperatorIpAddress     valueObject.IpAddress `json:"-"`
}

func NewDeleteSslPairVhosts(
	sslPairId valueObject.SslPairId,
	virtualHostsHostnames []valueObject.Fqdn,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteSslPairVhosts {
	return DeleteSslPairVhosts{
		SslPairId:             sslPairId,
		VirtualHostsHostnames: virtualHostsHostnames,
		OperatorAccountId:     operatorAccountId,
		OperatorIpAddress:     operatorIpAddress,
	}
}
