package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateSslPair struct {
	VirtualHostsHostnames []tkValueObject.Fqdn      `json:"virtualHostsHostnames"`
	Certificate           entity.SslCertificate     `json:"certificate"`
	ChainCertificates     *entity.SslCertificate    `json:"chainCertificates,omitempty"`
	Key                   valueObject.SslPrivateKey `json:"key"`
	OperatorAccountId     tkValueObject.AccountId   `json:"-"`
	OperatorIpAddress     tkValueObject.IpAddress   `json:"-"`
}

func NewCreateSslPair(
	virtualHostsHostnames []tkValueObject.Fqdn,
	certificate entity.SslCertificate,
	chainCertificates *entity.SslCertificate,
	key valueObject.SslPrivateKey,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
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
