package valueObject

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type SslPrivateKey string

func NewSslPrivateKey(privateKey string) (SslPrivateKey, error) {
	sslPrivateKey := SslPrivateKey(privateKey)
	if !sslPrivateKey.isValid() {
		return "", errors.New("SslPrivateKeyError")
	}

	return sslPrivateKey, nil
}

func NewSslPrivateKeyPanic(privateKey string) SslPrivateKey {
	sslPrivateKey, err := NewSslPrivateKey(privateKey)
	if err != nil {
		panic(err)
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
