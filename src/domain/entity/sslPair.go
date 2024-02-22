package entity

import "github.com/speedianet/os/src/domain/valueObject"

type SslPair struct {
	Id                valueObject.SslId         `json:"sslPairId"`
	VirtualHosts      []valueObject.Fqdn        `json:"virtualHosts"`
	Certificate       SslCertificate            `json:"certificate"`
	Key               valueObject.SslPrivateKey `json:"key"`
	ChainCertificates []SslCertificate          `json:"chainCertificates"`
}

func NewSslPair(
	sslPairId valueObject.SslId,
	virtualHosts []valueObject.Fqdn,
	certificate SslCertificate,
	key valueObject.SslPrivateKey,
	chainCertificates []SslCertificate,
) SslPair {
	return SslPair{
		Id:                sslPairId,
		VirtualHosts:      virtualHosts,
		Certificate:       certificate,
		Key:               key,
		ChainCertificates: chainCertificates,
	}
}
