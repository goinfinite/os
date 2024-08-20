package valueObject

import "testing"

func TestInlineHtmlContent(t *testing.T) {
	t.Run("ValidInlineHtmlContent", func(t *testing.T) {
		validInlineHtmlContents := []interface{}{
			"Some nice inline html content", "<h1>Nice title here</h1>",
			"<p>With some regular text here too...<h2>",
		}

		for _, inlineHtmlContent := range validInlineHtmlContents {
			_, err := NewInlineHtmlContent(inlineHtmlContent)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'",
					inlineHtmlContent, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidInlineHtmlContent", func(t *testing.T) {
		invalidInlineHtmlContents := []interface{}{
			"",
		}

		for _, inlineHtmlContent := range invalidInlineHtmlContents {
			_, err := NewInlineHtmlContent(inlineHtmlContent)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", inlineHtmlContent)
			}
		}
	})
}
