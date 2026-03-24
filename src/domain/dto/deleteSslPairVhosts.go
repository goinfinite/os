package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type DeleteSslPairVhosts struct {
	SslPairId             valueObject.SslPairId    `json:"sslPairId"`
	VirtualHostsHostnames []tkValueObject.Fqdn     `json:"virtualHostsHostnames"`
	OperatorAccountId     tkValueObject.AccountId  `json:"-"`
	OperatorIpAddress     tkValueObject.IpAddress  `json:"-"`
}

func NewDeleteSslPairVhosts(
	sslPairId valueObject.SslPairId,
	virtualHostsHostnames []tkValueObject.Fqdn,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteSslPairVhosts {
	return DeleteSslPairVhosts{
		SslPairId:             sslPairId,
		VirtualHostsHostnames: virtualHostsHostnames,
		OperatorAccountId:     operatorAccountId,
		OperatorIpAddress:     operatorIpAddress,
	}
}
