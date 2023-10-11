package dto

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type GetSslPair struct {
	SerialNumber      string           `json:"serialNumber"`
	Hostname          valueObject.Fqdn `json:"hostname"`
	Certificate       string           `json:"certificate"`
	Key               string           `json:"key"`
	ChainCertificates []string         `json:"chainCertificates"`
}

func NewGetSslPair(
	serialNumber valueObject.SslSerialNumber,
	hostname valueObject.Fqdn,
	certificate entity.SslCertificate,
	key entity.SslPrivateKey,
	chainCertificates []entity.SslCertificate,
) GetSslPair {
	var parsedChainCertificates []string
	for _, chainCertificate := range chainCertificates {
		parsedChainCertificates = append(parsedChainCertificates, chainCertificate.Certificate)
	}

	return GetSslPair{
		SerialNumber:      serialNumber.String(),
		Hostname:          hostname,
		Certificate:       certificate.Certificate,
		Key:               key.Key,
		ChainCertificates: parsedChainCertificates,
	}
}
