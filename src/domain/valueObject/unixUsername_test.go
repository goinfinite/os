package valueObject

import "testing"

func TestUnixUsername(t *testing.T) {
	t.Run("ValidUnixUsername", func(t *testing.T) {
		validUnixUsernames := []interface{}{
			"a", "a_1", "_abc-123", "b-c_d-e", "valid_name_with_30_chars",
		}

		for _, unixUsername := range validUnixUsernames {
			_, err := NewUnixUsername(unixUsername)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", unixUsername, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidUnixUsername", func(t *testing.T) {
		invalidUnixUsernames := []interface{}{
			"/1invalid_start_with_digit", "-invalid-start-with-dash",
			"invalid_character$more_than_30_chars",
			"toolongname_with_more_than_32_characters_long", "inv@lid_char",
		}

		for _, unixUsername := range invalidUnixUsernames {
			_, err := NewUnixUsername(unixUsername)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", unixUsername)
			}
		}
	})
}
