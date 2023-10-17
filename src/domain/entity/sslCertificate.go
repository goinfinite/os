package entity

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/speedianet/sam/src/domain/valueObject"
)

type SslCertificate struct {
	HashId      valueObject.SslHashId
	Certificate valueObject.SslCertificateStr
	CommonName  *valueObject.Fqdn
	IssuedAt    valueObject.UnixTime
	ExpiresAt   valueObject.UnixTime
	IsCA        bool
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

	certificate, err := valueObject.NewSslCertificateStr(sslCertificate)
	if err != nil {
		return SslCertificate{}, err
	}

	hashId, err := valueObject.NewSslHashIdFromSslCertificate(certificate)
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
		HashId:      hashId,
		Certificate: certificate,
		CommonName:  commonNamePtr,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiresAt,
		IsCA:        parsedCert.IsCA,
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
