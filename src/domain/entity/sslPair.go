package entity

import "github.com/speedianet/os/src/domain/valueObject"

type SslPair struct {
	Id                valueObject.SslId         `json:"sslPairId"`
	VirtualHost       valueObject.Fqdn          `json:"virtualHost"`
	Certificate       SslCertificate            `json:"certificate"`
	Key               valueObject.SslPrivateKey `json:"key"`
	ChainCertificates []SslCertificate          `json:"chainCertificates"`
}

func NewSslPair(
	sslPairId valueObject.SslId,
	virtualHost valueObject.Fqdn,
	certificate SslCertificate,
	key valueObject.SslPrivateKey,
	chainCertificates []SslCertificate,
) SslPair {
	return SslPair{
		Id:                sslPairId,
		VirtualHost:       virtualHost,
		Certificate:       certificate,
		Key:               key,
		ChainCertificates: chainCertificates,
	}
}
