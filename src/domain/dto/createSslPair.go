package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type CreateSslPair struct {
	VirtualHostsHostnames []valueObject.Fqdn        `json:"virtualHostsHostnames"`
	Certificate           entity.SslCertificate     `json:"certificate"`
	Key                   valueObject.SslPrivateKey `json:"key"`
	OperatorAccountId     valueObject.AccountId     `json:"-"`
	OperatorIpAddress     valueObject.IpAddress     `json:"-"`
}

func NewCreateSslPair(
	virtualHostsHostnames []valueObject.Fqdn,
	certificate entity.SslCertificate,
	key valueObject.SslPrivateKey,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateSslPair {
	return CreateSslPair{
		VirtualHostsHostnames: virtualHostsHostnames,
		Certificate:           certificate,
		Key:                   key,
		OperatorAccountId:     operatorAccountId,
		OperatorIpAddress:     operatorIpAddress,
	}
}
