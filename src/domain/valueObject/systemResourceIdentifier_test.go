package valueObject

import "testing"

func TestSystemResourceIdentifier(t *testing.T) {
	t.Run("ValidSystemResourceIdentifier", func(t *testing.T) {
		validSystemResourceIdentifier := []interface{}{
			"sri://0:account/120", "sri://1:mapping/200",
		}

		for _, identifier := range validSystemResourceIdentifier {
			_, err := NewSystemResourceIdentifier(identifier)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", identifier, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSystemResourceIdentifier", func(t *testing.T) {
		invalidSystemResourceIdentifier := []interface{}{
			"", "sri://0:/",
		}

		for _, identifier := range invalidSystemResourceIdentifier {
			_, err := NewSystemResourceIdentifier(identifier)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", identifier)
			}
		}
	})
}
