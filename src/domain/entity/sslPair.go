package entity

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

type SslPair struct {
	Certificate  string
	SerialNumber *big.Int
	CommonName   string
	IssuedAt     time.Time
	ExpiresAt    time.Time
}

func NewSslPair(certificate string) (SslPair, error) {
	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		return SslPair{}, errors.New("SslPairError")
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return SslPair{}, err
	}

	return SslPair{
		Certificate:  certificate,
		SerialNumber: parsedCert.SerialNumber,
		CommonName:   parsedCert.Subject.CommonName,
		IssuedAt:     parsedCert.NotBefore,
		ExpiresAt:    parsedCert.NotAfter,
	}, nil
}
