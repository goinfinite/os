package valueObject

import "testing"

func TestDataFieldValue(t *testing.T) {
	t.Run("ValidDataFieldValue", func(t *testing.T) {
		validDataFieldValues := []interface{}{
			"/",
			"This is my username",
			"new_email@mail.net",
			"localhost:8000",
			"https://www.google.com/search",
			1239218,
			1212.123,
			true,
			false,
		}

		for _, dfv := range validDataFieldValues {
			_, err := NewDataFieldValue(dfv)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", dfv, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldValue", func(t *testing.T) {
		invalidDataFieldValues := []interface{}{
			"",
		}

		for _, dfv := range invalidDataFieldValues {
			_, err := NewDataFieldValue(dfv)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dfv)
			}
		}
	})
}
