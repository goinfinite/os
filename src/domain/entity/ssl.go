package entity

import (
	"math/big"

	"github.com/speedianet/sam/src/domain/valueObject"
)

type Ssl struct {
	Id                big.Int          `json:"id"`
	Hostname          valueObject.Fqdn `json:"hostname"`
	Certificate       SslPair          `json:"certificate"`
	Key               SslPrivateKey    `json:"key"`
	ChainCertificates []SslPair        `json:"chainCertificates"`
}

func NewSsl(
	id big.Int,
	hostname valueObject.Fqdn,
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
