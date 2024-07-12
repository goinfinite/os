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
	IsIntermediary       bool                                `json:"-"`
	CertificateAuthority valueObject.SslCertificateAuthority `json:"certificateAuthority"`
	AltNames             []valueObject.SslHostname           `json:"altNames"`
	IssuedAt             valueObject.UnixTime                `json:"issuedAt"`
	ExpiresAt            valueObject.UnixTime                `json:"expiresAt"`
}

func NewSslCertificate(
	sslCertContent valueObject.SslCertificateContent,
) (SslCertificate, error) {
	var sslCertificate SslCertificate

	block, _ := pem.Decode([]byte(sslCertContent.String()))
	if block == nil {
		return sslCertificate, errors.New("SslCertificateContentDecodeError")
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return sslCertificate, errors.New("SslCertificateContentParseError")
	}

	sslCertId, err := valueObject.NewSslIdFromSslCertificateContent(sslCertContent)
	if err != nil {
		return sslCertificate, err
	}

	issuedAt := valueObject.NewUnixTimeWithGoTime(parsedCert.NotBefore)
	expiresAt := valueObject.NewUnixTimeWithGoTime(parsedCert.NotAfter)

	isIntermediary := true

	var commonNamePtr *valueObject.SslHostname
	commonName, err := valueObject.NewSslHostname(parsedCert.Subject.CommonName)
	if err == nil {
		isIntermediary = false
		commonNamePtr = &commonName
	}

	certAuthorityStr := "Self-signed"
	isSelfSigned := parsedCert.Subject.String() == parsedCert.Issuer.String()
	if !isSelfSigned {
		certIssuer := parsedCert.Issuer
		certAuthorityStr = certIssuer.CommonName

		hasOrganizationName := len(certIssuer.Organization) > 0 &&
			len(certIssuer.Organization[0]) > 0
		if hasOrganizationName {
			certAuthorityStr += ", " + certIssuer.Organization[0]
		}
	}

	certAuthority, err := valueObject.NewSslCertificateAuthority(certAuthorityStr)
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
		Id:                   sslCertId,
		CertificateContent:   sslCertContent,
		CommonName:           commonNamePtr,
		IsIntermediary:       isIntermediary,
		CertificateAuthority: certAuthority,
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
