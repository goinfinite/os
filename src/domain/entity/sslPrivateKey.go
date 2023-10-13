package entity

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type SslPrivateKey struct {
	Key string
}

func NewSslPrivateKey(privateKey string) (SslPrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return SslPrivateKey{}, errors.New("SslPrivateKeyError")
	}
	_, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return SslPrivateKey{}, nil
	}

	return SslPrivateKey{
		Key: privateKey,
	}, nil
}

func NewSslPrivateKeyPanic(privateKey string) SslPrivateKey {
	sslPrivateKey, err := NewSslPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	return sslPrivateKey
}

func (sslPrivateKey SslPrivateKey) String() string {
	return sslPrivateKey.Key
}
