package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type SslPair struct {
	SerialNumber      valueObject.SslSerialNumber `json:"serialNumber"`
	Hostname          valueObject.Fqdn            `json:"hostname"`
	Certificate       SslCertificate              `json:"certificate"`
	Key               valueObject.SslPrivateKey   `json:"key"`
	ChainCertificates []SslCertificate            `json:"chainCertificates"`
}

func NewSslPair(
	serialNumber valueObject.SslSerialNumber,
	hostname valueObject.Fqdn,
	certificate SslCertificate,
	key valueObject.SslPrivateKey,
	chainCertificates []SslCertificate,
) SslPair {
	return SslPair{
		SerialNumber:      serialNumber,
		Hostname:          hostname,
		Certificate:       certificate,
		Key:               key,
		ChainCertificates: chainCertificates,
	}
}
