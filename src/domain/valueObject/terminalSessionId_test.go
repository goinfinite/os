package valueObject

import "testing"

func TestTerminalSessionId(t *testing.T) {
	t.Run("ValidTerminalSessionId", func(t *testing.T) {
		validIds := []any{
			"a7f3c91d2e4b6f80", "0123456789abcdef", "1234567890123456",
			"A7F3C91D2E4B6F80",
		}

		for _, id := range validIds {
			_, err := NewTerminalSessionId(id)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", id, err.Error())
			}
		}
	})

	t.Run("InvalidTerminalSessionId", func(t *testing.T) {
		invalidIds := []any{
			"a7f3c91d2e4b6f8", "a7f3c91d2e4b6f800", "a7f3c91d2e4b6f8g",
			"infinite-a7f3c91d2e4b6f80", "", 12345,
		}

		for _, id := range invalidIds {
			_, err := NewTerminalSessionId(id)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", id)
			}
		}
	})
}
