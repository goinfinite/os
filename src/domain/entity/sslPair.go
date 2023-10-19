package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type SslPair struct {
	HashId            valueObject.SslId         `json:"hashId"`
	Hostname          valueObject.Fqdn          `json:"hostname"`
	Certificate       SslCertificate            `json:"certificate"`
	Key               valueObject.SslPrivateKey `json:"key"`
	ChainCertificates []SslCertificate          `json:"chainCertificates"`
}

func NewSslPair(
	hashId valueObject.SslId,
	hostname valueObject.Fqdn,
	certificate SslCertificate,
	key valueObject.SslPrivateKey,
	chainCertificates []SslCertificate,
) SslPair {
	return SslPair{
		HashId:            hashId,
		Hostname:          hostname,
		Certificate:       certificate,
		Key:               key,
		ChainCertificates: chainCertificates,
	}
}
