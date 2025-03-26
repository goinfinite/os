package entity

import (
	"time"

	"github.com/goinfinite/os/src/domain/valueObject"
)

const EarlyRenewalThresholdHours int64 = 48

type SslPair struct {
	Id                      valueObject.SslPairId     `json:"sslPairId"`
	MainVirtualHostHostname valueObject.Fqdn          `json:"mainVirtualHostHostname"`
	VirtualHostsHostnames   []valueObject.Fqdn        `json:"virtualHostsHostnames"`
	Certificate             SslCertificate            `json:"certificate"`
	Key                     valueObject.SslPrivateKey `json:"key"`
	ChainCertificates       []SslCertificate          `json:"chainCertificates"`
}

func NewSslPair(
	sslPairId valueObject.SslPairId,
	mainVirtualHostHostname valueObject.Fqdn,
	virtualHostsHostnames []valueObject.Fqdn,
	certificate SslCertificate,
	key valueObject.SslPrivateKey,
	chainCertificates []SslCertificate,
) SslPair {
	return SslPair{
		Id:                      sslPairId,
		MainVirtualHostHostname: mainVirtualHostHostname,
		VirtualHostsHostnames:   virtualHostsHostnames,
		Certificate:             certificate,
		Key:                     key,
		ChainCertificates:       chainCertificates,
	}
}

func (sslPair SslPair) IsPubliclyTrusted() bool {
	sslPairCrtAuthority := sslPair.Certificate.CertificateAuthority
	if sslPairCrtAuthority.IsSelfSigned() {
		return false
	}

	hoursToSeconds := int64(3600)
	earlyRenewalThresholdSeconds := EarlyRenewalThresholdHours * hoursToSeconds
	expirationDate := sslPair.Certificate.ExpiresAt.Int64()
	expirationDateUnixTime := time.Unix(expirationDate, 0).UTC().Unix()

	unixTimeToRenew := expirationDateUnixTime - earlyRenewalThresholdSeconds
	unixTimeNow := time.Now().Unix()
	shouldRenew := unixTimeNow >= unixTimeToRenew

	return !shouldRenew
}
