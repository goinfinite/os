package valueObject

import "testing"

func TestDataFieldType(t *testing.T) {
	t.Run("ValidDataFieldType", func(t *testing.T) {
		validDataFieldTypes := []interface{}{
			"checkbox", "color", "date", "email", "image", "number", "password",
			"radio", "range", "search", "tel", "text", "time", "url",
		}

		for _, dataFieldType := range validDataFieldTypes {
			_, err := NewDataFieldType(dataFieldType)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", dataFieldType, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidDataFieldType", func(t *testing.T) {
		invalidDataFieldTypes := []interface{}{
			"button", "datetime-local", "file", "hidden", "month", "reset",
			"submit", "week",
		}

		for _, dataFieldType := range invalidDataFieldTypes {
			_, err := NewDataFieldType(dataFieldType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", dataFieldType)
			}
		}
	})
}
