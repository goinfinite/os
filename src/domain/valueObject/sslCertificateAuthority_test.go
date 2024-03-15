package valueObject

import "testing"

func TestSslCertificateAuthority(t *testing.T) {
	t.Run("ValidSslCertificateAuthority", func(t *testing.T) {
		validsSslCertificateAuthority := []string{
			"IdenTrust",
			"DigiCert Group",
			"Sectigo (Comodo Cybersecurity)",
			"GlobalSign",
			"Let's Encrypt",
			"GoDaddy Group",
			"Internet Security Research Group",
		}
		for _, sslCertificateAuthority := range validsSslCertificateAuthority {
			_, err := NewSslCertificateAuthority(sslCertificateAuthority)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", sslCertificateAuthority, err.Error())
			}
		}
	})

	t.Run("InvalidSslCertificateAuthority", func(t *testing.T) {
		invalidsSslCertificateAuthority := []string{
			"",
			"Nitro Auth@rity",
			"Super long certificate authority, because I don't know",
		}
		for _, sslCertificateAuthority := range invalidsSslCertificateAuthority {
			_, err := NewSslCertificateAuthority(sslCertificateAuthority)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", sslCertificateAuthority)
			}
		}
	})
}
