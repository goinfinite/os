package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type AddSsl struct {
	Hostname    valueObject.VirtualHost    `json:"hostname"`
	Certificate valueObject.SslCertificate `json:"certificate"`
	Key         valueObject.SslPrivateKey  `json:"key"`
}

func NewAddSsl(
	hostname valueObject.VirtualHost,
	certificate valueObject.SslCertificate,
	key valueObject.SslPrivateKey,
) AddSsl {
	return AddSsl{
		Hostname:    hostname,
		Certificate: certificate,
		Key:         key,
	}
}
