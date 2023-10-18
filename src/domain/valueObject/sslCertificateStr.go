package valueObject

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type SslCertificateStr string

func NewSslCertificateStr(sslCertificate string) (SslCertificateStr, error) {
	certificate := SslCertificateStr(sslCertificate)
	if !certificate.isValid() {
		return "", errors.New("InvalidSslCertificateStr")
	}

	return certificate, nil
}

func NewSslCertificateStrPanic(certificate string) SslCertificateStr {
	sslCertificate, err := NewSslCertificateStr(certificate)
	if err != nil {
		panic(err)
	}
	return sslCertificate
}

func (sslCertificate SslCertificateStr) isValid() bool {
	block, _ := pem.Decode([]byte(sslCertificate))
	if block == nil {
		return false
	}

	_, err := x509.ParseCertificate(block.Bytes)
	return err != nil
}

func (sslCertificate SslCertificateStr) String() string {
	return string(sslCertificate)
}
