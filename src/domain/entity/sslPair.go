package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type SslPair struct {
	sslPairId         valueObject.SslId         `json:"sslPairId"`
	Hostname          valueObject.Fqdn          `json:"hostname"`
	Certificate       SslCertificate            `json:"certificate"`
	Key               valueObject.SslPrivateKey `json:"key"`
	ChainCertificates []SslCertificate          `json:"chainCertificates"`
}

func NewSslPair(
	sslPairId valueObject.SslId,
	hostname valueObject.Fqdn,
	certificate SslCertificate,
	key valueObject.SslPrivateKey,
	chainCertificates []SslCertificate,
) SslPair {
	return SslPair{
		sslPairId:         sslPairId,
		Hostname:          hostname,
		Certificate:       certificate,
		Key:               key,
		ChainCertificates: chainCertificates,
	}
}
