package entity

import (
	"testing"
	"time"

	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

func TestSslPair(t *testing.T) {
	t.Run("SelfSignedSslPair", func(t *testing.T) {
		testHostname, err := tkValueObject.NewFqdn("test.example.com")
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %s", err.Error())
		}

		synthesizer := &tkInfra.Synthesizer{}
		certContentStr, keyContentStr, err := synthesizer.SelfSignedCertificatePairPemFactory(
			&testHostname, nil,
		)
		if err != nil {
			t.Fatalf("SelfSignedCertificatePairGenerationFailed: %s", err.Error())
		}

		certificateContent, err := valueObject.NewSslCertificateContent(certContentStr)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %s", err.Error())
		}
		certificate, err := NewSslCertificate(certificateContent)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %s", err.Error())
		}

		privateKey, err := valueObject.NewSslPrivateKey(keyContentStr)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %s", err.Error())
		}

		pairId, err := valueObject.NewSslPairIdFromSslPairContent(
			certificateContent, []valueObject.SslCertificateContent{}, privateKey,
		)
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %s", err.Error())
		}

		sslPair := NewSslPair(
			pairId, testHostname, certificate, privateKey, []SslCertificate{},
		)

		if sslPair.IsPubliclyTrusted() {
			t.Errorf("SelfSignedSslPairShouldNotBePubliclyTrusted")
		}
	})

	t.Run("PubliclyTrustedSslPair", func(t *testing.T) {
		certificateAuthority, err := valueObject.NewSslCertificateAuthority("Test CA")
		if err != nil {
			t.Fatalf("ExpectedNoErrorButGot: %s", err.Error())
		}

		sslPair := SslPair{
			Certificate: SslCertificate{
				CertificateAuthority: certificateAuthority,
				ExpiresAt:            tkValueObject.NewUnixTimeAfterNow(30 * 24 * time.Hour),
			},
		}

		if !sslPair.IsPubliclyTrusted() {
			t.Errorf("PubliclyTrustedSslPairShouldBePubliclyTrusted")
		}
	})
}
