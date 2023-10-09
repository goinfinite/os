package dto

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type GetSsl struct {
	Id                string           `json:"id"`
	Hostname          valueObject.Fqdn `json:"hostname"`
	Certificate       string           `json:"certificate"`
	Key               string           `json:"key"`
	ChainCertificates []string         `json:"chainCertificates"`
}

func NewGetSsl(
	id valueObject.SslId,
	hostname valueObject.Fqdn,
	certificate entity.SslPair,
	key entity.SslPrivateKey,
	chainCertificates []entity.SslPair,
) GetSsl {
	var parsedChainCertificates []string
	for _, chainCertificate := range chainCertificates {
		parsedChainCertificates = append(parsedChainCertificates, chainCertificate.Certificate)
	}

	return GetSsl{
		Id:                id.String(),
		Hostname:          hostname,
		Certificate:       certificate.Certificate,
		Key:               key.Key,
		ChainCertificates: parsedChainCertificates,
	}
}
