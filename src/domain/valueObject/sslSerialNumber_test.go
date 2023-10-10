package valueObject

import (
	"testing"
)

func TestNewSslSerialNumber(t *testing.T) {
	t.Run("ValidSerialNumber", func(t *testing.T) {
		validSerialNumbers := []string{
			"12345",
			"314159265358979323846264338327950288419716939937510582097494459",
			"598437582340",
			"21094819052730572857384801928390218492359374608430598239047894",
		}

		for _, sslSerialNumber := range validSerialNumbers {
			_, err := NewSslSerialNumber(sslSerialNumber)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", sslSerialNumber, err)
			}
		}
	})

	t.Run("InvalidSerialNumber", func(t *testing.T) {
		invalidSerialNumbers := []string{
			"-1",
			"0",
			"-1231231231242353467",
		}

		for _, sslSerialNumber := range invalidSerialNumbers {
			_, err := NewSslSerialNumber(sslSerialNumber)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", sslSerialNumber)
			}
		}
	})
}
