package entity

import (
	"time"

	"github.com/speedianet/os/src/domain/valueObject"
)

const EarlyRenewalThresholdHours int64 = 48

type SslPair struct {
	Id                    valueObject.SslId         `json:"sslPairId"`
	VirtualHostsHostnames []valueObject.Fqdn        `json:"virtualHostsHostnames"`
	Certificate           SslCertificate            `json:"certificate"`
	Key                   valueObject.SslPrivateKey `json:"key"`
	ChainCertificates     []SslCertificate          `json:"chainCertificates"`
}

func NewSslPair(
	sslPairId valueObject.SslId,
	virtualHostsHostnames []valueObject.Fqdn,
	certificate SslCertificate,
	key valueObject.SslPrivateKey,
	chainCertificates []SslCertificate,
) SslPair {
	return SslPair{
		Id:                    sslPairId,
		VirtualHostsHostnames: virtualHostsHostnames,
		Certificate:           certificate,
		Key:                   key,
		ChainCertificates:     chainCertificates,
	}
}

func (sslPair SslPair) IsPubliclyTrusted() bool {
	sslPairCrtAuthority := sslPair.Certificate.CertificateAuthority
	if sslPairCrtAuthority.IsSelfSigned() {
		return false
	}

	hoursToSeconds := int64(3600)
	earlyRenewalThresholdSeconds := EarlyRenewalThresholdHours * hoursToSeconds
	expirationDate := sslPair.Certificate.ExpiresAt.Read()
	expirationDateUnixTime := time.Unix(expirationDate, 0).UTC().Unix()

	unixTimeToRenew := expirationDateUnixTime - earlyRenewalThresholdSeconds
	unixTimeNow := time.Now().Unix()
	shouldRenew := unixTimeNow >= unixTimeToRenew

	return !shouldRenew
}
