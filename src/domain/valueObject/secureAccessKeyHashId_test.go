package valueObject

import (
	"testing"
)

func TestSecureAccessKeyHashId(t *testing.T) {
	t.Run("ValidSecureAccessKeyHashId", func(t *testing.T) {
		rawValidSecureAccessKeyHashId := []interface{}{
			"q09qgpsmyqQg3QolSYgQSUYp", "rEJqbCBzJWsSKysWYkHmpsVa",
			"kuLH8x2t96AIrPwcBr0kW2fH", "mW5f4ZQzUgxi2kGzZT+aC5Tf",
		}

		for _, rawKeyHashId := range rawValidSecureAccessKeyHashId {
			_, err := NewSecureAccessKeyHashId(rawKeyHashId)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", rawKeyHashId, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSecureAccessKeyHashId", func(t *testing.T) {
		rawInvalidSecureAccessKeyHashId := []interface{}{
			"", 1.50, true, 1000,
		}

		for _, rawKeyHashId := range rawInvalidSecureAccessKeyHashId {
			_, err := NewSecureAccessKeyHashId(rawKeyHashId)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", rawKeyHashId)
			}
		}
	})
}
