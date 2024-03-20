package entity

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/speedianet/os/src/domain/valueObject"
)

type SslCertificate struct {
	Id                   valueObject.SslId                   `json:"sslId"`
	CommonName           *valueObject.SslHostname            `json:"commonName"`
	CertificateContent   valueObject.SslCertificateContent   `json:"certificateContent"`
	IsCA                 bool                                `json:"-"`
	CertificateAuthority valueObject.SslCertificateAuthority `json:"certificateAuthority"`
	IssuerCommonName     valueObject.SslHostname             `json:"-"`
	AltNames             []valueObject.SslHostname           `json:"altNames"`
	IssuedAt             valueObject.UnixTime                `json:"issuedAt"`
	ExpiresAt            valueObject.UnixTime                `json:"expiresAt"`
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
	if !parsedCert.IsCA {
		commonName, err := valueObject.NewSslHostname(parsedCert.Subject.CommonName)
		if err != nil {
			return sslCertificate, errors.New("InvalidSslCertificateCommonName")
		}
		commonNamePtr = &commonName
	}

	certIssuer := parsedCert.Issuer
	issuerCommonNameStr := certIssuer.CommonName
	issuerCommonName, err := valueObject.NewSslHostname(issuerCommonNameStr)
	if err != nil {
		return sslCertificate, errors.New("InvalidIssuerCommonName")
	}

	if len(certIssuer.Organization) == 0 {
		return sslCertificate, errors.New("SslCertificateWithoutCA")
	}

	certificateAuthorityStr := certIssuer.Organization[0]
	certificateAuthority, err := valueObject.NewSslCertificateAuthority(certificateAuthorityStr)
	if err != nil {
		return sslCertificate, err
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
		Id:                   sslCertificateId,
		CertificateContent:   sslCertificateContent,
		CommonName:           commonNamePtr,
		IsCA:                 parsedCert.IsCA,
		CertificateAuthority: certificateAuthority,
		IssuerCommonName:     issuerCommonName,
		AltNames:             altNames,
		IssuedAt:             issuedAt,
		ExpiresAt:            expiresAt,
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
