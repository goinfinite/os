package valueObject

import (
	"testing"
)

func TestNewSslId(t *testing.T) {
	t.Run("ValidSslId", func(t *testing.T) {
		validSslIds := []interface{}{
			"a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4",
			"a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4",
			"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			"ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890",
			"0f1e2d3c4b5a697887a6b5c4d3e2f1a01234567890abcdef1234567890abcdef",
		}

		for _, validSslId := range validSslIds {
			_, err := NewSslId(validSslId)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", validSslId, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSslId", func(t *testing.T) {
		invalidSslIds := []interface{}{
			"g3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3d4e5f6a7b8c9d0e1f2a3b4",
			"12345", "!@#$%^&*()_+|}{:?><,./;'[]=-",
			"abcdefgh1234567890abcdefgh1234567890abcdefgh1234567890abcdefgh12",
			"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcde",
		}

		for _, invalidSslId := range invalidSslIds {
			_, err := NewSslId(invalidSslId)
			if err == nil {
				t.Errorf("Expected no error for '%v', got nil", invalidSslId)
			}
		}
	})

	t.Run("ValidSslIdFromSslPairContent", func(t *testing.T) {
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
		validChainCertsContent := []SslCertificateContent{validCertContent}

		validKey := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDXGqUB2HVTK+rz
UtbpytCYiO6zeMHNe22n4jgXTsJHoj/BvwE5S5kZirIkjCIAyAX+Pu6blLnCnftJ
Gb8X90K5GqZRGFc3vX3oDj/Umtlk1zkBkyqVDjDBy/BxLm0ddKHgjTZHxzOHaC/5
etnNGbw1y23MWDwH1SvRBqxpp5bISbfh80sO22FLU1hDb5Pr4f5uUeCMKU/n0G5V
YN0UwP55RpwqjhD7nPb6l01g7r6MBY4EVAR76EM1kNcpd1czijMY+4LNI0ZQj2+l
bFQ02yx2Yuw/8wAPeuOPaToYMYe0er7YTET0L8rQEytMcKskYvO9Q5Ekrn3l9JYw
t+eB+QjDAgMBAAECggEBAI2HCmpccVV33+6Y4q6Qsw6pieSr31fDjjKXtTAgsdNP
/YMMmVGJXAJiLzO8v+KjuM2/ul7DTDWwnFVMi17JYS1JS4Sv7zLNirnUJktMVxzy
Pp+6pJnN7GaWOG0/jquCwb9tKfmwJ0dAVdBf9E3uUNdUMbnxlA3TRDETov0hNyQv
p1ah3Z00QbmnZfzeVpv26m+nKysAH6yj+qCOc5yAZNvTzIE5AW4KFykuIZCpVjmX
s0Q8NBQ+sPVQZnQHzba1dtfCOakDQnOn3UtYJWly5F8Cu+seZPElZlhrp/b7jI0S
rABn5K994IYwQN3AE+rucmZ1uET5vcWevM+JVxL6LgECgYEA+ct9PwXzYQtbbsA5
yU/ZYCRsi2dmGeRF5Ltb8v05vj1z/5BT4HPw4Gx8kYv8jaOH0M/O5qoYjhaYZA7g
rKRAIXGfQmzncY1GYOZWZ959w+CwCr1bOv4TXSSoEJLIXjYtXJsNBlVk5Ne+M1T1
P8NMxcgyW7u5OckFFbIgL+P8spECgYEA3HKMJ5Zua1EbiSqnMJGBH6CEkVrXGH2t
c3EeiYQshtWPsFSpUB6767Hib3dpFBvMaPj2xWFrvxI70ru0Ag4VRMZ2poTrHay2
Ge+wR1Q5zMnmxqP35mmW4YargIRYU+ctORcJo73W3fd47Z/lIlxrBabvx4mHJ89A
/VKTTvWPSBMCgYBJT4Fol7SADLc+38eV34tqfgYlO6lpe+dPY/VucQcbYCnFHXSg
cSaGlxBQHwd2AkJ/9B1C8TTXrqX2567kvCfeKNyWwCOE3fODyNYfEdtTO4QvArfd
rme8dF+mzY1kqP3TKeY+r5021GKL6ik4F3dWrJSq+4M3BForreVoaL7nIQKBgBLE
IZBBKxcxqWFs4xysVkyl8oMZM5RfJoPcTlgwi0XTKk89dchfRWoUE42foa9Xingp
MYCuAWkbmUIgPnuqTT80kectC4LUMBBXKi94SQ9Y5K9mR/Uyaei6+SCQo6BI3r2s
a2KoB4GPzpiT8wKQ0X+CrYjT+VB3QTYPcIDZQKBHAoGAC8gDLUfHWDA+Ozuj7fZT
9NYN1ALwBoHC10bTWDAw9dC+l5p2yv0qJ8waaJrbXevuQGbH/+WItsZnVt+CxjfC
5jXdBpt0nixwIinr970lG2kQc2Jf64VtS9KoRoO2qnHVfNcn0DnVoWTvRjjeqVxx
PZIyej7kPh0NXWwDyV9uhyk=
-----END PRIVATE KEY-----`
		validKeyContent, _ := NewSslPrivateKey(validKey)

		_, err := NewSslIdFromSslPairContent(
			validCertContent, validChainCertsContent, validKeyContent,
		)
		if err != nil {
			t.Errorf(
				"Expected no error for '%v', got '%s'", validCertContent, err.Error(),
			)
		}
	})

	t.Run("ValidSslIdFromSslCertificateContent", func(t *testing.T) {
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
		_, err := NewSslIdFromSslCertificateContent(validCertContent)
		if err != nil {
			t.Errorf(
				"Expected no error for '%v', got '%s'", validCertContent, err.Error(),
			)
		}
	})
}
