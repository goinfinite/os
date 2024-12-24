package valueObject

import "testing"

func TestDataFieldSpecificType(t *testing.T) {
	t.Run("ValidDataFieldSpecificType", func(t *testing.T) {
		validDataFieldSpecificTypes := []interface{}{
			"password", "PASSWORD", "username", "USERNAME", "email", "EMAIL",
		}

		for _, dataFieldSpecificType := range validDataFieldSpecificTypes {
			_, err := NewDataFieldSpecificType(dataFieldSpecificType)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", dataFieldSpecificType,
					err.Error(),
				)
			}
		}
	})

	t.Run("InvalidDataFieldSpecificType", func(t *testing.T) {
		invalidDataFieldSpecificTypes := []interface{}{
			"button", "datetime-local", "file", "hidden", "month", "reset",
			"submit", "week",
		}

		for _, dataFieldSpecificType := range invalidDataFieldSpecificTypes {
			_, err := NewDataFieldSpecificType(dataFieldSpecificType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", dataFieldSpecificType)
			}
		}
	})
}
