package valueObject

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestInlineHtmlContent(t *testing.T) {
	t.Run("ValidInlineHtmlContent", func(t *testing.T) {
		validInlineHtmlContents := []string{
			"Some nice inline html content",
			"<h1>Nice title here</h1>",
			"<p>With some regular text here too...<h2>",
		}

		for _, ihc := range validInlineHtmlContents {
			_, err := NewInlineHtmlContent(ihc)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", ihc, err.Error())
			}
		}
	})

	t.Run("ValidInlineHtmlContent", func(t *testing.T) {
		invalidLength := 3600
		invalidInlineHtmlContents := []string{
			"",
			testHelpers.GenerateString(invalidLength),
		}

		for _, ihc := range invalidInlineHtmlContents {
			_, err := NewInlineHtmlContent(ihc)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", ihc)
			}
		}
	})
}
