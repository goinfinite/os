package valueObject

import (
	"testing"
)

func TestNewSslCertificateId(t *testing.T) {
	t.Run("ValidSslCertificateId", func(t *testing.T) {
		validSslCertificateIds := []interface{}{
			"a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4",
			"a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4",
			"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			"ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890",
			"0f1e2d3c4b5a697887a6b5c4d3e2f1a01234567890abcdef1234567890abcdef",
		}

		for _, validSslCertificateId := range validSslCertificateIds {
			_, err := NewSslCertificateId(validSslCertificateId)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", validSslCertificateId, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSslCertificateId", func(t *testing.T) {
		invalidSslCertificateIds := []interface{}{
			"g3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4",
			"12345", "!@#$%^&*()_+|}{:?><,./;'[]=-",
			"abcdefgh1234567890abcdefgh1234567890abcdefgh1234567890abcdefgh12",
			"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde",
		}

		for _, invalidSslCertificateId := range invalidSslCertificateIds {
			_, err := NewSslCertificateId(invalidSslCertificateId)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", invalidSslCertificateId)
			}
		}
	})

	t.Run("ValidSslCertificateIdFromSslCertificateContent", func(t *testing.T) {
		validCert := `
-----BEGIN CERTIFICATE-----
MIIDujCCAqKgAwIBAgIIE31FZVaPXTUwDQYJKoZIhvcNAQEFBQAwSTELMAkGA1UE
BhMCVVMxEzARBgNVBAoTCkdvb2dsZSBJbmMxJTAjBgNVBAMTHEdvb2dsZSBJbnRl
cm5ldCBBdXRob3JpdHkgRzIwHhcNMTQwMTI5MTMyNzQzWhcNMTQwNTI5MDAwMDAw
WjBpMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwN
TW91bnRhaW4gVmlldzETMBEGA1UECgwKR29vZ2xlIEluYzEYMBYGA1UEAwwPbWFp
bC5nb29nbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEfRrObuSW5T7q
5CnSEqefEmtH4CCv6+5EckuriNr1CjfVvqzwfAhopXkLrq45EQm8vkmf7W96XJhC
7ZM0dYi1/qOCAU8wggFLMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAa
BgNVHREEEzARgg9tYWlsLmdvb2dsZS5jb20wCwYDVR0PBAQDAgeAMGgGCCsGAQUF
BwEBBFwwWjArBggrBgEFBQcwAoYfaHR0cDovL3BraS5nb29nbGUuY29tL0dJQUcy
LmNydDArBggrBgEFBQcwAYYfaHR0cDovL2NsaWVudHMxLmdvb2dsZS5jb20vb2Nz
cDAdBgNVHQ4EFgQUiJxtimAuTfwb+aUtBn5UYKreKvMwDAYDVR0TAQH/BAIwADAf
BgNVHSMEGDAWgBRK3QYWG7z2aLV29YG2u2IaulqBLzAXBgNVHSAEEDAOMAwGCisG
AQQB1nkCBQEwMAYDVR0fBCkwJzAloCOgIYYfaHR0cDovL3BraS5nb29nbGUuY29t
L0dJQUcyLmNybDANBgkqhkiG9w0BAQUFAAOCAQEAH6RYHxHdcGpMpFE3oxDoFnP+
gtuBCHan2yE2GRbJ2Cw8Lw0MmuKqHlf9RSeYfd3BXeKkj1qO6TVKwCh+0HdZk283
TZZyzmEOyclm3UGFYe82P/iDFt+CeQ3NpmBg+GoaVCuWAARJN/KfglbLyyYygcQq
0SgeDh8dRKUiaW3HQSoYvTvdTuqzwK4CXsr3b5/dAOY8uMuG/IAR3FgwTbZ1dtoW
RvOTa8hYiU6A475WuZKyEHcwnGYe57u2I2KbMgcKjPniocj4QzgYsVAVKW3IwaOh
yE+vPxsiUkvQHdO2fojCkY8jg70jxM+gu59tPDNbw3Uh/2Ij310FgTHsnGQMyA==
-----END CERTIFICATE-----`

		validCertContent, _ := NewSslCertificateContent(validCert)
		_, err := NewSslCertificateIdFromSslCertificateContent(validCertContent)
		if err != nil {
			t.Errorf(
				"Expected no error for '%v', got '%s'", validCertContent, err.Error(),
			)
		}
	})
}
