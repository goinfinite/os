package dto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type AddSslPair struct {
	VirtualHost valueObject.Fqdn          `json:"virtualHost"`
	Certificate entity.SslCertificate     `json:"certificate"`
	Key         valueObject.SslPrivateKey `json:"key"`
}

func NewAddSslPair(
	virtualHost valueObject.Fqdn,
	certificate entity.SslCertificate,
	key valueObject.SslPrivateKey,
) AddSslPair {
	return AddSslPair{
		VirtualHost: virtualHost,
		Certificate: certificate,
		Key:         key,
	}
}
