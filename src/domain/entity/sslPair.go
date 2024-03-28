package entity

import (
	"time"

	"github.com/speedianet/os/src/domain/valueObject"
)

const EarlyRenewalThresholdHours int64 = 48

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

func (sslPair SslPair) IsPubliclyTrusted() bool {
	sslPairCrtAuthority := sslPair.Certificate.CertificateAuthority
	if sslPairCrtAuthority.IsSelfSigned() {
		return false
	}

	hoursToSeconds := int64(3600)
	earlyRenewalThresholdSeconds := EarlyRenewalThresholdHours * hoursToSeconds

	expirationDate := sslPair.Certificate.ExpiresAt.Get()
	expirationDateUnixTime := time.Unix(expirationDate, 0).UTC().Unix()
	unixTimeToRenew := expirationDateUnixTime - earlyRenewalThresholdSeconds
	nowUnixTime := time.Now().Unix()

	return nowUnixTime > unixTimeToRenew
}
