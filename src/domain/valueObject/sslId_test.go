package valueObject

import (
	"testing"
)

func TestNewSslId(t *testing.T) {
	t.Run("ValidId", func(t *testing.T) {
		validIds := []string{
			"12345",
			"314159265358979323846264338327950288419716939937510582097494459",
			"598437582340",
			"21094819052730572857384801928390218492359374608430598239047894",
		}

		for _, sslId := range validIds {
			_, err := NewSslId(sslId)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", sslId, err)
			}
		}
	})

	t.Run("InvalidId", func(t *testing.T) {
		invalidIds := []string{
			"-1",
			"0",
			"-1231231231242353467",
		}

		for _, sslId := range invalidIds {
			_, err := NewSslId(sslId)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", sslId)
			}
		}
	})
}
