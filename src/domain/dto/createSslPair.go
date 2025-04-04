package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type CreateSslPair struct {
	VirtualHostsHostnames []valueObject.Fqdn        `json:"virtualHostsHostnames"`
	Certificate           entity.SslCertificate     `json:"certificate"`
	ChainCertificates     *entity.SslCertificate    `json:"chainCertificates,omitempty"`
	Key                   valueObject.SslPrivateKey `json:"key"`
	OperatorAccountId     valueObject.AccountId     `json:"-"`
	OperatorIpAddress     valueObject.IpAddress     `json:"-"`
}

func NewCreateSslPair(
	virtualHostsHostnames []valueObject.Fqdn,
	certificate entity.SslCertificate,
	chainCertificates *entity.SslCertificate,
	key valueObject.SslPrivateKey,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateSslPair {
	return CreateSslPair{
		VirtualHostsHostnames: virtualHostsHostnames,
		Certificate:           certificate,
		ChainCertificates:     chainCertificates,
		Key:                   key,
		OperatorAccountId:     operatorAccountId,
		OperatorIpAddress:     operatorIpAddress,
	}
}
