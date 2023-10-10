package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type SslPair struct {
	Id                valueObject.SslSerialNumber `json:"id"`
	Hostname          valueObject.Fqdn            `json:"hostname"`
	Certificate       SslCertificate              `json:"certificate"`
	Key               SslPrivateKey               `json:"key"`
	ChainCertificates []SslCertificate            `json:"chainCertificates"`
}

func NewSslPair(
	id valueObject.SslSerialNumber,
	hostname valueObject.Fqdn,
	certificate SslCertificate,
	key SslPrivateKey,
	chainCertificates []SslCertificate,
) SslPair {
	return SslPair{
		Id:                id,
		Hostname:          hostname,
		Certificate:       certificate,
		Key:               key,
		ChainCertificates: chainCertificates,
	}
}
