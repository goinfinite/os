package valueObject

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type SslCertificateContent string

func NewSslCertificateContent(input interface{}) (
	certContent SslCertificateContent, err error,
) {
	stringValue, err := voHelper.InterfaceToString(input)
	if err != nil {
		return certContent, errors.New("SslCertificateContentMustBeString")
	}

	pemBlock, _ := pem.Decode([]byte(stringValue))
	if pemBlock == nil {
		return certContent, errors.New("InvalidSslCertificateContent")
	}

	_, err = x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return certContent, errors.New("InvalidSslCertificateContent")
	}

	return SslCertificateContent(stringValue), nil
}

func (vo SslCertificateContent) String() string {
	return string(vo)
}

func (vo SslCertificateContent) MarshalJSON() ([]byte, error) {
	voBytes := []byte(string(vo))
	return json.Marshal(base64.StdEncoding.EncodeToString(voBytes))
}
