package dto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type AddSslPair struct {
	VirtualHosts []valueObject.Fqdn        `json:"virtualHosts"`
	Certificate  entity.SslCertificate     `json:"certificate"`
	Key          valueObject.SslPrivateKey `json:"key"`
}

func NewAddSslPair(
	virtualHosts []valueObject.Fqdn,
	certificate entity.SslCertificate,
	key valueObject.SslPrivateKey,
) AddSslPair {
	return AddSslPair{
		VirtualHosts: virtualHosts,
		Certificate:  certificate,
		Key:          key,
	}
}
