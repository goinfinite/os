package valueObject

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type SslPrivateKey string

func NewSslPrivateKey(value string) (SslPrivateKey, error) {
	sslPrivateKey := SslPrivateKey(value)
	if !sslPrivateKey.isValid() {
		return "", errors.New("InvalidSslPrivateKey")
	}
	return sslPrivateKey, nil
}

func NewSslPrivateKeyPanic(value string) SslPrivateKey {
	sslPrivateKey := SslPrivateKey(value)
	if !sslPrivateKey.isValid() {
		panic("InvalidSslPrivateKey")
	}
	return sslPrivateKey
}

func (sslPrivateKey SslPrivateKey) isValid() bool {
	block, _ := pem.Decode([]byte(sslPrivateKey))
	if block == nil {
		return false
	}
	_, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return false
	}

	return true
}

func (sslPrivateKey SslPrivateKey) String() string {
	return string(sslPrivateKey)
}
