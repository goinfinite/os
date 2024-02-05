package dto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type CreateSslPair struct {
	Hostname    valueObject.Fqdn          `json:"hostname"`
	Certificate entity.SslCertificate     `json:"certificate"`
	Key         valueObject.SslPrivateKey `json:"key"`
}

func NewCreateSslPair(
	hostname valueObject.Fqdn,
	certificate entity.SslCertificate,
	key valueObject.SslPrivateKey,
) CreateSslPair {
	return CreateSslPair{
		Hostname:    hostname,
		Certificate: certificate,
		Key:         key,
	}
}
