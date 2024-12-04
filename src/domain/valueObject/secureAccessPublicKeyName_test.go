package valueObject

import (
	"testing"
)

func TestSecureAccessPublicKeyName(t *testing.T) {
	t.Run("ValidSecureAccessPublicKeyName", func(t *testing.T) {
		rawValidSecureAccessPublicKeyName := []interface{}{
			"myMachine@pop-os", "thats-my-only-pc", "tryingWithThisTypeOfName",
		}

		for _, rawKeyName := range rawValidSecureAccessPublicKeyName {
			_, err := NewSecureAccessPublicKeyName(rawKeyName)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", rawKeyName, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSecureAccessPublicKeyName", func(t *testing.T) {
		rawInvalidSecureAccessPublicKeyName := []interface{}{
			"", "that's not allowed, u know?", "maybe-with-#",
			"thisIsAnEnormousNameToTestVoLength",
		}

		for _, rawKeyName := range rawInvalidSecureAccessPublicKeyName {
			_, err := NewSecureAccessPublicKeyName(rawKeyName)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", rawKeyName)
			}
		}
	})
}
