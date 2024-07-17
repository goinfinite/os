package valueObject

import (
	"testing"
)

func TestUsername(t *testing.T) {
	t.Run("ValidUsername", func(t *testing.T) {
		validUsernames := []string{
			"a",
			"a_1",
			"_abc-123",
			"b-c_d-e",
			"valid_name_with_30_chars",
		}

		for _, username := range validUsernames {
			_, err := NewUsername(username)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", username, err)
			}
		}
	})

	t.Run("InvalidUsername", func(t *testing.T) {
		invalidUsernames := []string{
			"/1invalid_start_with_digit",
			"-invalid-start-with-dash",
			"invalid_character$more_than_30_chars",
			"toolongname_with_more_than_32_characters_long",
			"inv@lid_char",
		}

		for _, username := range invalidUsernames {
			_, err := NewUsername(username)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", username)
			}
		}
	})
}
