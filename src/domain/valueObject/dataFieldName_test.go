package valueObject

import (
	"testing"
)

func TestDataFieldName(t *testing.T) {
	t.Run("ValidDataFieldName", func(t *testing.T) {
		validDataFieldNames := []string{
			"username",
			"user-email",
			"Service-Name_With_Port80",
		}

		for _, dfn := range validDataFieldNames {
			_, err := NewDataFieldName(dfn)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", dfn, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldName", func(t *testing.T) {
		invalidDataFieldNames := []string{
			"",
			"./test",
			"-key",
			"anotherkey-",
		}

		for _, dfn := range invalidDataFieldNames {
			_, err := NewDataFieldName(dfn)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dfn)
			}
		}
	})
}
