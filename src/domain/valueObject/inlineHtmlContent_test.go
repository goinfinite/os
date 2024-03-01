package valueObject

import (
	"testing"
)

func TestNewInlineHtmlContent(t *testing.T) {
	t.Run("ValidInlineHtml", func(t *testing.T) {
		validInlineHtmls := []string{
			"<html><p>This is a HTML <strong>test</strong>.</p></html>",
			"<html><div class='container'><span>Content</span></div></html>",
			"<html><img src='image.jpg' alt='image'/></html>",
			"<html>Test without tags</html>",
		}

		for _, inlineHtml := range validInlineHtmls {
			_, err := NewInlineHtmlContent(inlineHtml)
			if err != nil {
				t.Errorf("ExpectingNoErrorButGot: %s", err.Error())
			}
		}
	})

	t.Run("InvalidInlineHtml", func(t *testing.T) {
		invalidInlineHtmls := []string{
			"",
			"12345",
			"Texto sem tags HTML",
			"<p>HTML com erro",
		}

		for _, inlineHtml := range invalidInlineHtmls {
			_, err := NewInlineHtmlContent(inlineHtml)
			if err == nil {
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", inlineHtml)
			}
		}
	})
}
