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
		return "", errors.New("InvalidSslPrivateKey")
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
	_, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	return err == nil
}

func NewSslPrivateKeyFromEncodedContent(
	encodedContent EncodedContent,
) (SslPrivateKey, error) {
	var sslPrivateKey SslPrivateKey

	decodedContent, err := encodedContent.GetDecodedContent()
	if err != nil {
		return sslPrivateKey, errors.New("InvalidSslPrivateKey")
	}

	return NewSslPrivateKey(decodedContent)
}

func NewSslPrivateKeyFromEncodedContentPanic(
	encodedContent EncodedContent,
) SslPrivateKey {
	decodedContent, err := encodedContent.GetDecodedContent()
	if err != nil {
		panic("InvalidSslPrivateKey")
	}

	return NewSslPrivateKeyPanic(decodedContent)
}

func (sslPrivateKey SslPrivateKey) String() string {
	return string(sslPrivateKey)
}
