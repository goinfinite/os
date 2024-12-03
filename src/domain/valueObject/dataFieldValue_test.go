package valueObject

import "testing"

func TestDataFieldValue(t *testing.T) {
	t.Run("ValidDataFieldValue", func(t *testing.T) {
		validDataFieldValues := []interface{}{
			"/", "This is my username", "new_email@mail.net", "localhost:8000",
			"https://www.google.com/search", 1239218, 1212.123, true, false,
		}

		for _, value := range validDataFieldValues {
			_, err := NewDataFieldValue(value)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", value, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldValue", func(t *testing.T) {
		invalidDataFieldValues := []interface{}{
			"",
		}

		for _, value := range invalidDataFieldValues {
			_, err := NewDataFieldValue(value)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", value)
			}
		}
	})
}
