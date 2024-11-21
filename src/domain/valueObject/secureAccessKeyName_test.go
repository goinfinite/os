package valueObject

import (
	"testing"
)

func TestSecureAccessKeyName(t *testing.T) {
	t.Run("ValidSecureAccessKeyName", func(t *testing.T) {
		rawValidSecureAccessKeyName := []interface{}{
			"myMachine@pop-os", "thats-my-only-pc", "tryingWithThisTypeOfName",
		}

		for _, rawKeyName := range rawValidSecureAccessKeyName {
			_, err := NewSecureAccessKeyName(rawKeyName)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", rawKeyName, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSecureAccessKeyName", func(t *testing.T) {
		rawInvalidSecureAccessKeyName := []interface{}{
			"", "that's not allowed, u know?", "maybe-with-#",
			"thisIsAnEnormousNameToTestVoLength",
		}

		for _, rawKeyName := range rawInvalidSecureAccessKeyName {
			_, err := NewSecureAccessKeyName(rawKeyName)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", rawKeyName)
			}
		}
	})
}
