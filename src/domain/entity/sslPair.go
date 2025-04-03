package entity

import (
	"time"

	"github.com/goinfinite/os/src/domain/valueObject"
)

const EarlyRenewalThresholdHours int64 = 48

type SslPair struct {
	Id                  valueObject.SslPairId     `json:"sslPairId"`
	VirtualHostHostname valueObject.Fqdn          `json:"virtualHostHostname"`
	Certificate         SslCertificate            `json:"certificate"`
	Key                 valueObject.SslPrivateKey `json:"key"`
	ChainCertificates   []SslCertificate          `json:"chainCertificates"`
}

func NewSslPair(
	sslPairId valueObject.SslPairId,
	virtualHostHostname valueObject.Fqdn,
	certificate SslCertificate,
	key valueObject.SslPrivateKey,
	chainCertificates []SslCertificate,
) SslPair {
	return SslPair{
		Id:                  sslPairId,
		VirtualHostHostname: virtualHostHostname,
		Certificate:         certificate,
		Key:                 key,
		ChainCertificates:   chainCertificates,
	}
}

func (sslPair SslPair) IsPubliclyTrusted() bool {
	sslPairCrtAuthority := sslPair.Certificate.CertificateAuthority
	if sslPairCrtAuthority.IsSelfSigned() {
		return false
	}

	expirationDate := sslPair.Certificate.ExpiresAt.ReadAsGoTime().UTC()
	earlyRenewalThreshold := time.Hour * time.Duration(EarlyRenewalThresholdHours)
	renewalDeadline := expirationDate.Add(-earlyRenewalThreshold)

	currentTime := time.Now().UTC()
	return currentTime.Before(renewalDeadline)
}
