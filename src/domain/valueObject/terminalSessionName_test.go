package valueObject

import (
	"strings"
	"testing"
)

func TestTerminalSessionName(t *testing.T) {
	t.Run("ValidTerminalSessionName", func(t *testing.T) {
		validNames := []any{
			"my session", "  my session  ", "Web-1_prod (blue).2", "a",
			"1session", strings.Repeat("a", 64),
		}

		for _, name := range validNames {
			_, err := NewTerminalSessionName(name)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", name, err.Error())
			}
		}
	})

	t.Run("InvalidTerminalSessionName", func(t *testing.T) {
		invalidNames := []any{
			"", "   ", "-session", "a|b", "a/b", "sessão",
			strings.Repeat("a", 65), "a\nb", []string{"a"},
		}

		for _, name := range invalidNames {
			_, err := NewTerminalSessionName(name)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", name)
			}
		}
	})
}
