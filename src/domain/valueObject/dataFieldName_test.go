package valueObject

import "testing"

func TestDataFieldName(t *testing.T) {
	t.Run("ValidDataFieldName", func(t *testing.T) {
		validDataFieldNames := []interface{}{
			"username", "user-email", "Service-Name_With_Port80",
		}

		for _, name := range validDataFieldNames {
			_, err := NewDataFieldName(name)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", name, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldName", func(t *testing.T) {
		invalidDataFieldNames := []interface{}{
			"", "./test", "-key", "anotherkey-",
		}

		for _, name := range invalidDataFieldNames {
			_, err := NewDataFieldName(name)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", name)
			}
		}
	})
}
