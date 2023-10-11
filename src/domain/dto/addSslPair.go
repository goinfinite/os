package dto

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type AddSslPair struct {
	Hostname    valueObject.Fqdn      `json:"hostname"`
	Certificate entity.SslCertificate `json:"certificate"`
	Key         entity.SslPrivateKey  `json:"key"`
}

func NewAddSslPair(
	hostname valueObject.Fqdn,
	certificate entity.SslCertificate,
	key entity.SslPrivateKey,
) AddSslPair {
	return AddSslPair{
		Hostname:    hostname,
		Certificate: certificate,
		Key:         key,
	}
}
