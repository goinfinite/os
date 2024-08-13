package valueObject

import (
	"testing"
)

func TestFqdn(t *testing.T) {
	t.Run("ValidFqdn", func(t *testing.T) {
		validFqdns := []string{
			"example.com",
			"sub.example.com",
			"*.example.com",
			"my-site.co.uk",
			"sub-domain.example.org",
		}

		for _, fqdn := range validFqdns {
			_, err := NewFqdn(fqdn)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", fqdn, err.Error())
			}
		}
	})

	t.Run("InvalidFqdn", func(t *testing.T) {
		invalidFqdns := []string{
			"-example.com",
			"example-.com",
			"example.c",
			"example..com",
			"*example.com",
		}

		for _, fqdn := range invalidFqdns {
			_, err := NewFqdn(fqdn)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", fqdn)
			}
		}
	})
}
