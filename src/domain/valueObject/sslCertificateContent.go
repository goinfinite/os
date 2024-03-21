package valueObject

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
)

type SslCertificateContent string

func NewSslCertificateContent(sslCertificate string) (SslCertificateContent, error) {
	certificate := SslCertificateContent(sslCertificate)
	if !certificate.isValid() {
		return "", errors.New("InvalidSslCertificateContent")
	}

	return certificate, nil
}

func NewSslCertificateContentPanic(certificate string) SslCertificateContent {
	sslCertificate, err := NewSslCertificateContent(certificate)
	if err != nil {
		panic(err)
	}
	return sslCertificate
}

func (sslCrt SslCertificateContent) isValid() bool {
	block, _ := pem.Decode([]byte(sslCrt))
	if block == nil {
		return false
	}

	_, err := x509.ParseCertificate(block.Bytes)
	return err == nil
}

func NewSslCertificateContentFromEncodedContent(
	encodedContent EncodedContent,
) (SslCertificateContent, error) {
	var sslCertificateContent SslCertificateContent

	decodedContent, err := encodedContent.GetDecodedContent()
	if err != nil {
		return sslCertificateContent, errors.New("InvalidSslCertificate")
	}

	return NewSslCertificateContent(decodedContent)
}

func NewSslCertificateContentFromEncodedContentPanic(
	encodedContent EncodedContent,
) SslCertificateContent {
	decodedContent, err := encodedContent.GetDecodedContent()
	if err != nil {
		panic("InvalidSslCertificate")
	}

	return NewSslCertificateContentPanic(decodedContent)
}

func (sslCrt SslCertificateContent) String() string {
	return string(sslCrt)
}

func (sslCrt SslCertificateContent) MarshalJSON() ([]byte, error) {
	sslCrtBytes := []byte(string(sslCrt))
	return json.Marshal(base64.StdEncoding.EncodeToString(sslCrtBytes))
}
