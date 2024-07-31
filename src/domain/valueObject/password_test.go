package valueObject

import "testing"

func TestPassword(t *testing.T) {
	t.Run("ValidPassword", func(t *testing.T) {
		validPasswords := []interface{}{
			"password123", "S3cureP@ssw0rd!",
			"A_longer_password_with_various_chars123!", "MySecret2024",
			"Th1s!s@G00dPass",
		}

		for _, password := range validPasswords {
			_, err := NewPassword(password)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", password, err.Error())
			}
		}
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		invalidPasswords := []interface{}{
			"short", "tiny", "abc", "pass", "p@ss1",
		}

		for _, password := range invalidPasswords {
			_, err := NewPassword(password)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", password)
			}
		}
	})
}
