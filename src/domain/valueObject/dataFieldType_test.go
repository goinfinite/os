package valueObject

import "testing"

func TestDataFieldType(t *testing.T) {
	t.Run("ValidDataFieldType", func(t *testing.T) {
		validDataFieldTypes := []string{
			"checkbox",
			"color",
			"date",
			"email",
			"image",
			"number",
			"password",
			"radio",
			"range",
			"search",
			"tel",
			"text",
			"time",
			"url",
		}

		for _, dft := range validDataFieldTypes {
			_, err := NewDataFieldType(dft)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", dft, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldType", func(t *testing.T) {
		invalidDataFieldTypes := []string{
			"button",
			"datetime-local",
			"file",
			"hidden",
			"month",
			"reset",
			"submit",
			"week",
		}

		for _, dft := range invalidDataFieldTypes {
			_, err := NewDataFieldType(dft)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dft)
			}
		}
	})
}
