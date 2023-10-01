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
	sslCertificate := SslCertificate(value)
	if !sslCertificate.isValid() {
		panic("InvalidSslCertificate")
	}
	return sslCertificate
}

func (sslCertificate SslCertificate) isValid() bool {
	block, _ := pem.Decode([]byte(sslCertificate))
	if block == nil {
		return false
	}
	_, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false
	}

	return true
}

func (sslCertificate SslCertificate) String() string {
	return string(sslCertificate)
}
