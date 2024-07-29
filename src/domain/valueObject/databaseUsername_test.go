package valueObject

import (
	"testing"
)

func TestDatabaseUsername(t *testing.T) {
	t.Run("ValidDatabaseUsername", func(t *testing.T) {
		validDatabaseUsernames := []string{
			"abc",
			"a_1",
			"_abc-123",
			"b-c_d-e",
			"valid_name_with_30_chars",
		}

		for _, dbUsername := range validDatabaseUsernames {
			_, err := NewDatabaseUsername(dbUsername)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", dbUsername, err)
			}
		}
	})

	t.Run("InvalidDatabaseUsername", func(t *testing.T) {
		invalidDatabaseUsernames := []string{
			"/1invalid_start_with_digit",
			"-invalid-start-with-dash",
			"invalid_character$more_than_30_chars",
			"inv@lid_char",
		}

		for _, dbUsername := range invalidDatabaseUsernames {
			_, err := NewDatabaseUsername(dbUsername)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dbUsername)
			}
		}
	})
}
