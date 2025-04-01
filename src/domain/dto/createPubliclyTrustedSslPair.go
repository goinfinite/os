package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type CreatePubliclyTrustedSslPair struct {
	CommonName        valueObject.Fqdn      `json:"commonName"`
	AltNames          []valueObject.Fqdn    `json:"aliasesHostnames"`
	OperatorAccountId valueObject.AccountId `json:"-"`
	OperatorIpAddress valueObject.IpAddress `json:"-"`
}

func NewCreatePubliclyTrustedSslPair(
	commonName valueObject.Fqdn,
	altNames []valueObject.Fqdn,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreatePubliclyTrustedSslPair {
	return CreatePubliclyTrustedSslPair{
		CommonName:        commonName,
		AltNames:          altNames,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
