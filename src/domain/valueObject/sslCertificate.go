package valueObject

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type SslCertificate string

func NewSslCertificate(value string) (SslCertificate, error) {
	sslCertificate := SslCertificate(value)
	if !sslCertificate.isValid() {
		return "", errors.New("InvalidSslCertificate")
	}
	return sslCertificate, nil
}

func NewSslCertificatePanic(value string) SslCertificate {
	sslCertificate, err := NewSslCertificate(value)
	if err != nil {
		panic(err)
	}
	return sslCertificate
}

func (sslCertificate SslCertificate) isValid() bool {
	_, err := sslCertificate.GetCertInfo()
	if err != nil {
		return false
	}

	return true
}

func (sslCertificate SslCertificate) GetCertInfo() (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(sslCertificate))
	if block == nil {
		return nil, errors.New("PemDecodeError")
	}
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return parsedCert, nil
}

func (sslCertificate SslCertificate) String() string {
	return string(sslCertificate)
}
