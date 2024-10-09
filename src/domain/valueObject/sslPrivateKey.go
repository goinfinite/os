package valueObject

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
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

func (vo SslPrivateKey) String() string {
	return string(vo)
}

func (vo SslPrivateKey) MarshalJSON() ([]byte, error) {
	voBytes := []byte(string(vo))
	return json.Marshal(base64.StdEncoding.EncodeToString(voBytes))
}
