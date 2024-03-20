package valueObject

import (
	"crypto/x509"
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

func (sslCertificate SslCertificateContent) isValid() bool {
	block, _ := pem.Decode([]byte(sslCertificate))
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

func (sslCertificate SslCertificateContent) String() string {
	return string(sslCertificate)
}
