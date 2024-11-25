package valueObject

import (
	"testing"
)

func TestSecureAccessKeyUuid(t *testing.T) {
	t.Run("ValidSecureAccessKeyUuid", func(t *testing.T) {
		rawValidSecureAccessKeyUuid := []interface{}{
			"abc123def4", "1234567890ab", "abcdef123456", "9876543210ab",
			"1234abcd5678",
		}

		for _, rawKeyUuid := range rawValidSecureAccessKeyUuid {
			_, err := NewSecureAccessKeyUuid(rawKeyUuid)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", rawKeyUuid, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSecureAccessKeyUuid", func(t *testing.T) {
		rawInvalidSecureAccessKeyUuid := []interface{}{
			1234, true, 1.40, "abc123", "tooLongSecureAccessKeyUuid", "12345678!@#",
			"short12",
		}

		for _, rawKeyUuid := range rawInvalidSecureAccessKeyUuid {
			_, err := NewSecureAccessKeyUuid(rawKeyUuid)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", rawKeyUuid)
			}
		}
	})
}
