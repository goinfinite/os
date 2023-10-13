package entity

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

type SslCertificate struct {
	Certificate  string
	SerialNumber *big.Int
	CommonName   string
	IssuedAt     time.Time
	ExpiresAt    time.Time
	IsCA         bool
}

func NewSslCertificate(certificate string) (SslCertificate, error) {
	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		return SslCertificate{}, errors.New("SslCertificateError")
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return SslCertificate{}, err
	}

	return SslCertificate{
		Certificate:  certificate,
		SerialNumber: parsedCert.SerialNumber,
		CommonName:   parsedCert.Subject.CommonName,
		IssuedAt:     parsedCert.NotBefore,
		ExpiresAt:    parsedCert.NotAfter,
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
	return sslCertificate.Certificate
}
