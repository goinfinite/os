package valueObject

import (
	"testing"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

func TestNewSslPairId(t *testing.T) {
	t.Run("ValidSslPairId", func(t *testing.T) {
		validSslPairIds := []any{
			"a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4",
			"a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4",
			"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			"ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890",
			"0f1e2d3c4b5a697887a6b5c4d3e2f1a01234567890abcdef1234567890abcdef",
		}

		for _, validSslPairId := range validSslPairIds {
			_, err := NewSslPairId(validSslPairId)
			if err != nil {
				t.Errorf(
					"UnexpectedError: %v, value: '%v'", err.Error(), validSslPairId,
				)
			}
		}
	})

	t.Run("InvalidSslPairId", func(t *testing.T) {
		invalidSslPairIds := []any{
			"g3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4",
			"12345", "!@#$%^&*()_+|}{:?><,./;'[]=-",
			"abcdefgh1234567890abcdefgh1234567890abcdefgh1234567890abcdefgh12",
			"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde",
		}

		for _, invalidSslPairId := range invalidSslPairIds {
			_, err := NewSslPairId(invalidSslPairId)
			if err == nil {
				t.Errorf("MissingExpectedError: '%v' was accepted", invalidSslPairId)
			}
		}
	})

	t.Run("ValidSslPairIdFromSslPairContent", func(t *testing.T) {
		hostname, err := tkValueObject.NewFqdn("pair-id.example.com")
		if err != nil {
			t.Fatalf("NewFqdnFailed: %s", err.Error())
		}

		synthesizer := &tkInfra.Synthesizer{}
		certContentStr, keyContentStr, err := synthesizer.SelfSignedCertificatePairPemFactory(
			&hostname, nil,
		)
		if err != nil {
			t.Fatalf("SelfSignedCertificatePairGenerationFailed: %s", err.Error())
		}

		certContent, err := NewSslCertificateContent(certContentStr)
		if err != nil {
			t.Fatalf("SslCertificateContentParseFailed: %v", err)
		}
		keyContent, err := NewSslPrivateKey(keyContentStr)
		if err != nil {
			t.Fatalf("SslPrivateKeyParseFailed: %v", err)
		}
		chainCertsContent := []SslCertificateContent{certContent}

		sslPairId, err := NewSslPairIdFromSslPairContent(
			certContent, chainCertsContent, keyContent,
		)
		if err != nil {
			t.Fatalf("SslPairIdBuildFailed: %v", err)
		}

		repeatedSslPairId, err := NewSslPairIdFromSslPairContent(
			certContent, chainCertsContent, keyContent,
		)
		if err != nil {
			t.Fatalf("SslPairIdBuildFailed: %v", err)
		}
		if sslPairId != repeatedSslPairId {
			t.Error("SslPairIdShouldBeDeterministic")
		}

		otherKeyContentStr, err := synthesizer.PrivateKeyPemFactory(
			tkInfra.PrivateKeySettings{
				Algorithm: tkValueObject.PrivateKeyAlgorithmECDSA,
				BitSize:   256,
			},
		)
		if err != nil {
			t.Fatalf("FailedToGenerateSslPrivateKey: %s", err.Error())
		}
		otherPrivateKey, err := NewSslPrivateKey(otherKeyContentStr)
		if err != nil {
			t.Fatalf("SslPrivateKeyParseFailed: %v", err)
		}

		otherSslPairId, err := NewSslPairIdFromSslPairContent(
			certContent, chainCertsContent, otherPrivateKey,
		)
		if err != nil {
			t.Fatalf("SslPairIdBuildFailed: %v", err)
		}
		if sslPairId == otherSslPairId {
			t.Error("SslPairIdShouldChangeWithPrivateKey")
		}
	})

	t.Run("ChainCertificateCountChangesPairId", func(t *testing.T) {
		certContent := SslCertificateContent("chain-edge-fixture-cert")
		keyContent := SslPrivateKey("chain-edge-fixture-key")

		emptyChainSslPairId, err := NewSslPairIdFromSslPairContent(
			certContent, []SslCertificateContent{}, keyContent,
		)
		if err != nil {
			t.Fatalf("SslPairIdBuildFailed: %v", err)
		}

		singleChainSslPairId, err := NewSslPairIdFromSslPairContent(
			certContent, []SslCertificateContent{certContent}, keyContent,
		)
		if err != nil {
			t.Fatalf("SslPairIdBuildFailed: %v", err)
		}

		if emptyChainSslPairId == singleChainSslPairId {
			t.Error("SslPairIdShouldChangeWithChainCount")
		}
	})
}
