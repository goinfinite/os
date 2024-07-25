package valueObject

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type SslPrivateKey string

func NewSslPrivateKey(value interface{}) (privateKey SslPrivateKey, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return privateKey, errors.New("SslPrivateKeyMustBeString")
	}

	pemBlock, _ := pem.Decode([]byte(stringValue))
	if pemBlock == nil {
		return privateKey, errors.New("InvalidSslPrivateKey")
	}

	_, err = x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		_, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
		if err != nil {
			return privateKey, errors.New("InvalidSslPrivateKey")
		}
	}

	return SslPrivateKey(stringValue), nil
}

func NewSslPrivateKeyPanic(privateKey string) SslPrivateKey {
	sslPrivateKey, err := NewSslPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	return sslPrivateKey
}

func NewSslPrivateKeyFromEncodedContent(
	encodedContent EncodedContent,
) (privateKey SslPrivateKey, err error) {
	decodedContent, err := encodedContent.GetDecodedContent()
	if err != nil {
		return privateKey, errors.New("InvalidSslPrivateKey")
	}

	return NewSslPrivateKey(decodedContent)
}

// TODO: Remover isso.
func NewSslPrivateKeyFromEncodedContentPanic(
	encodedContent EncodedContent,
) SslPrivateKey {
	decodedContent, err := encodedContent.GetDecodedContent()
	if err != nil {
		panic("InvalidSslPrivateKey")
	}

	return NewSslPrivateKeyPanic(decodedContent)
}

func (vo SslPrivateKey) String() string {
	return string(vo)
}

func (vo SslPrivateKey) MarshalJSON() ([]byte, error) {
	voBytes := []byte(string(vo))
	return json.Marshal(base64.StdEncoding.EncodeToString(voBytes))
}
