package dto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type CreateSslPair struct {
	VirtualHosts []valueObject.Fqdn        `json:"virtualHosts"`
	Certificate  entity.SslCertificate     `json:"certificate"`
	Key          valueObject.SslPrivateKey `json:"key"`
}

func NewCreateSslPair(
	virtualHosts []valueObject.Fqdn,
	certificate entity.SslCertificate,
	key valueObject.SslPrivateKey,
) CreateSslPair {
	return CreateSslPair{
		VirtualHosts: virtualHosts,
		Certificate:  certificate,
		Key:          key,
	}
}
