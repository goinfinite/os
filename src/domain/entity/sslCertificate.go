package entity

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/speedianet/sam/src/domain/valueObject"
)

type SslCertificate struct {
	Certificate  valueObject.SslCertificateStr
	SerialNumber valueObject.SslSerialNumber
	CommonName   *valueObject.Fqdn
	IssuedAt     valueObject.UnixTime
	ExpiresAt    valueObject.UnixTime
	IsCA         bool
}

func NewSslCertificate(sslCertificate string) (SslCertificate, error) {
	block, _ := pem.Decode([]byte(sslCertificate))
	if block == nil {
		return SslCertificate{}, errors.New("SslCertificateError")
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return SslCertificate{}, err
	}

	certificate := valueObject.NewSslCertificateStrPanic(sslCertificate)
	serialNumber := valueObject.NewSslSerialNumberPanic(parsedCert.SerialNumber)
	issuedAt := valueObject.UnixTime(parsedCert.NotBefore.Unix())
	expiresAt := valueObject.UnixTime(parsedCert.NotAfter.Unix())

	var commonNamePtr *valueObject.Fqdn
	commonNamePtr = nil
	if !parsedCert.IsCA {
		commonName := valueObject.NewFqdnPanic(parsedCert.Subject.CommonName)
		commonNamePtr = &commonName
	}

	return SslCertificate{
		Certificate:  certificate,
		SerialNumber: serialNumber,
		CommonName:   commonNamePtr,
		IssuedAt:     issuedAt,
		ExpiresAt:    expiresAt,
		IsCA:         parsedCert.IsCA,
	}, nil
}

func NewSslCertificatePanic(certificate string) SslCertificate {
	sslCertificate, err := NewSslCertificate(certificate)
	if err != nil {
		panic(err)
	}
	return sslCertificate
}

func (sslCertificate SslCertificate) String() string {
	return sslCertificate.Certificate.String()
}
