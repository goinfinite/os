package entity

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/goinfinite/os/src/domain/valueObject"
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
	certContent valueObject.SslCertificateContent,
) (certificate SslCertificate, err error) {
	block, _ := pem.Decode([]byte(certContent.String()))
	if block == nil {
		return certificate, errors.New("SslCertificateContentDecodeError")
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return certificate, errors.New("SslCertificateContentParseError")
	}

	certId, err := valueObject.NewSslIdFromSslCertificateContent(certContent)
	if err != nil {
		return certificate, err
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
		return certificate, err
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
		Id:                   certId,
		CertificateContent:   certContent,
		CommonName:           commonNamePtr,
		IsIntermediary:       isIntermediary,
		CertificateAuthority: certAuthority,
		AltNames:             altNames,
		IssuedAt:             issuedAt,
		ExpiresAt:            expiresAt,
	}, nil
}
