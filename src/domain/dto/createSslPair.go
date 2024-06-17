package dto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type CreateSslPair struct {
	VirtualHostsHostnames []valueObject.Fqdn        `json:"virtualHostsHostnames"`
	Certificate           entity.SslCertificate     `json:"certificate"`
	Key                   valueObject.SslPrivateKey `json:"key"`
}

func NewCreateSslPair(
	virtualHostsHostnames []valueObject.Fqdn,
	certificate entity.SslCertificate,
	key valueObject.SslPrivateKey,
) CreateSslPair {
	return CreateSslPair{
		VirtualHostsHostnames: virtualHostsHostnames,
		Certificate:           certificate,
		Key:                   key,
	}
}
