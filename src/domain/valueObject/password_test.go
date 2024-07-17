package valueObject

import (
	"testing"
)

func TestPassword(t *testing.T) {
	t.Run("ValidPassword", func(t *testing.T) {
		validPassword := []string{
			"password123",
			"S3cureP@ssw0rd!",
			"A_longer_password_with_various_chars123!",
			"MySecret2024",
			"Th1s!s@G00dPass",
		}

		for _, password := range validPassword {
			_, err := NewPassword(password)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", password, err)
			}
		}
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		invalidPassword := []string{
			"short",
			"tiny",
			"abc",
			"pass",
			"p@ss1",
		}

		for _, password := range invalidPassword {
			_, err := NewPassword(password)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", password)
			}
		}
	})
}
