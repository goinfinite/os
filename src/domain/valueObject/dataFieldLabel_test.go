package valueObject

import (
	"testing"
)

func TestDataFieldLabel(t *testing.T) {
	t.Run("ValidDataFieldLabel", func(t *testing.T) {
		validDataFieldLabels := []string{
			"checkbock",
			"color",
			"date",
			"datetime-local",
			"email",
			"image",
			"month",
			"number",
			"password",
			"radio",
			"range",
			"search",
			"tel",
			"text",
			"time",
			"url",
			"week",
		}

		for _, dfl := range validDataFieldLabels {
			_, err := NewDataFieldLabel(dfl)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", dfl, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldLabel", func(t *testing.T) {
		invalidDataFieldLabels := []string{
			"button",
			"file",
			"hidden",
			"reset",
			"submit",
		}

		for _, dfl := range invalidDataFieldLabels {
			_, err := NewDataFieldLabel(dfl)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dfl)
			}
		}
	})
}
