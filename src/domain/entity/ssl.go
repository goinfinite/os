package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type Ssl struct {
	Id                valueObject.SslId       `json:"id"`
	Hostname          valueObject.VirtualHost `json:"hostname"`
	Certificate       SslPair                 `json:"certificate"`
	Key               SslPrivateKey           `json:"key"`
	ChainCertificates []SslPair               `json:"chainCertificates"`
}

func NewSsl(
	id valueObject.SslId,
	hostname valueObject.VirtualHost,
	certificate SslPair,
	key SslPrivateKey,
	chainCertificates []SslPair,
) Ssl {
	return Ssl{
		Id:                id,
		Hostname:          hostname,
		Certificate:       certificate,
		Key:               key,
		ChainCertificates: chainCertificates,
	}
}
