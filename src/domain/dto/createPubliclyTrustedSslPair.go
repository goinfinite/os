package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type CreatePubliclyTrustedSslPair struct {
	VirtualHostHostname valueObject.Fqdn      `json:"virtualHostHostname"`
	OperatorAccountId   valueObject.AccountId `json:"-"`
	OperatorIpAddress   valueObject.IpAddress `json:"-"`
}

func NewCreatePubliclyTrustedSslPair(
	virtualHostHostname valueObject.Fqdn,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreatePubliclyTrustedSslPair {
	return CreatePubliclyTrustedSslPair{
		VirtualHostHostname: virtualHostHostname,
		OperatorAccountId:   operatorAccountId,
		OperatorIpAddress:   operatorIpAddress,
	}
}
