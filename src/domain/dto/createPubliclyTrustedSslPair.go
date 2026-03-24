package dto

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type CreatePubliclyTrustedSslPair struct {
	VirtualHostHostname tkValueObject.Fqdn      `json:"virtualHostHostname"`
	OperatorAccountId   tkValueObject.AccountId `json:"-"`
	OperatorIpAddress   tkValueObject.IpAddress `json:"-"`
}

func NewCreatePubliclyTrustedSslPair(
	virtualHostHostname tkValueObject.Fqdn,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreatePubliclyTrustedSslPair {
	return CreatePubliclyTrustedSslPair{
		VirtualHostHostname: virtualHostHostname,
		OperatorAccountId:   operatorAccountId,
		OperatorIpAddress:   operatorIpAddress,
	}
}
