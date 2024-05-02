package valueObject

import (
	"testing"
)

func TestHtmlType(t *testing.T) {
	t.Run("ValidHtmlType", func(t *testing.T) {
		validHtmlTypes := []string{
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

		for _, ht := range validHtmlTypes {
			_, err := NewHtmlType(ht)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", ht, err.Error())
			}
		}
	})

	t.Run("InvalidHtmlType", func(t *testing.T) {
		invalidHtmlTypes := []string{
			"button",
			"file",
			"hidden",
			"reset",
			"submit",
		}

		for _, ht := range invalidHtmlTypes {
			_, err := NewHtmlType(ht)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", ht)
			}
		}
	})
}
