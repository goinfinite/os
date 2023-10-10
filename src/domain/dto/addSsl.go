package dto

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type AddSsl struct {
	Hostname    valueObject.Fqdn      `json:"hostname"`
	Certificate entity.SslCertificate `json:"certificate"`
	Key         entity.SslPrivateKey  `json:"key"`
}

func NewAddSsl(
	hostname valueObject.Fqdn,
	certificate entity.SslCertificate,
	key entity.SslPrivateKey,
) AddSsl {
	return AddSsl{
		Hostname:    hostname,
		Certificate: certificate,
		Key:         key,
	}
}
