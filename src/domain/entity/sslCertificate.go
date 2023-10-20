package entity

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/speedianet/sam/src/domain/valueObject"
)

type SslCertificate struct {
	SslCertificateId valueObject.SslId
	Certificate      valueObject.SslCertificateContent
	CommonName       *valueObject.Fqdn
	IssuedAt         valueObject.UnixTime
	ExpiresAt        valueObject.UnixTime
	IsCA             bool
}

func NewSslCertificate(
	sslCertificateContent valueObject.SslCertificateContent,
) (SslCertificate, error) {
	block, _ := pem.Decode([]byte(sslCertificateContent.String()))
	if block == nil {
		return SslCertificate{}, errors.New("SslCertificateError")
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return SslCertificate{}, err
	}

	sslCertificateId, err := valueObject.NewSslIdFromSslCertificateContent(
		sslCertificateContent,
	)
	if err != nil {
		return SslCertificate{}, err
	}

	issuedAt := valueObject.UnixTime(parsedCert.NotBefore.Unix())
	expiresAt := valueObject.UnixTime(parsedCert.NotAfter.Unix())

	var commonNamePtr *valueObject.Fqdn
	commonNamePtr = nil
	if !parsedCert.IsCA {
		commonName, err := valueObject.NewFqdn(parsedCert.Subject.CommonName)
		if err != nil {
			return SslCertificate{}, err
		}
		commonNamePtr = &commonName
	}

	return SslCertificate{
		SslCertificateId: sslCertificateId,
		Certificate:      sslCertificateContent,
		CommonName:       commonNamePtr,
		IssuedAt:         issuedAt,
		ExpiresAt:        expiresAt,
		IsCA:             parsedCert.IsCA,
	}, nil
}

func NewSslCertificatePanic(
	sslCertificateContent valueObject.SslCertificateContent,
) SslCertificate {
	sslCertificate, err := NewSslCertificate(sslCertificateContent)
	if err != nil {
		panic(err)
	}
	return sslCertificate
}

func (sslCertificate SslCertificate) String() string {
	return sslCertificate.Certificate.String()
}
