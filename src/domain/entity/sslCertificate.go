package entity

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/speedianet/os/src/domain/valueObject"
)

type SslCertificate struct {
	Id                 valueObject.SslId                 `json:"sslId"`
	CertificateContent valueObject.SslCertificateContent `json:"certificateContent"`
	CommonName         *valueObject.SslHostname          `json:"commonName"`
	IsCA               bool                              `json:"isCa"`
	AltNames           []valueObject.SslHostname         `json:"altNames"`
	IssuedAt           valueObject.UnixTime              `json:"issuedAt"`
	ExpiresAt          valueObject.UnixTime              `json:"expiresAt"`
}

func NewSslCertificate(
	sslCertificateContent valueObject.SslCertificateContent,
) (SslCertificate, error) {
	var sslCertificate SslCertificate

	block, _ := pem.Decode([]byte(sslCertificateContent.String()))
	if block == nil {
		return sslCertificate, errors.New("SslCertificateContentDecodeError")
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return sslCertificate, errors.New("SslCertificateContentParseError")
	}

	sslCertificateId, err := valueObject.NewSslIdFromSslCertificateContent(
		sslCertificateContent,
	)
	if err != nil {
		return sslCertificate, err
	}

	issuedAt := valueObject.UnixTime(parsedCert.NotBefore.Unix())
	expiresAt := valueObject.UnixTime(parsedCert.NotAfter.Unix())

	var commonNamePtr *valueObject.SslHostname
	commonNamePtr = nil
	if !parsedCert.IsCA {
		commonName, err := valueObject.NewSslHostname(parsedCert.Subject.CommonName)
		if err != nil {
			return sslCertificate, errors.New("InvalidSslCertificateCommonName")
		}
		commonNamePtr = &commonName
	}

	altNames := []valueObject.SslHostname{}
	if len(parsedCert.DNSNames) > 0 {
		for _, certDnsName := range parsedCert.DNSNames {
			altName, err := valueObject.NewSslHostname(certDnsName)
			if err != nil {
				continue
			}

			altNames = append(altNames, altName)
		}
	}

	return SslCertificate{
		Id:                 sslCertificateId,
		CertificateContent: sslCertificateContent,
		CommonName:         commonNamePtr,
		IsCA:               parsedCert.IsCA,
		AltNames:           altNames,
		IssuedAt:           issuedAt,
		ExpiresAt:          expiresAt,
	}, nil
}

func NewSslCertificatePanic(
	sslCertificateContent valueObject.SslCertificateContent,
) SslCertificate {
	sslCertificate, err := NewSslCertificate(sslCertificateContent)
	if err != nil {
		panic(err)
	}
	return sslCertificate
}
