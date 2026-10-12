package valueObject

import "testing"

func TestUnixCompressionType(t *testing.T) {
	t.Run("ValidUnixCompressionType", func(t *testing.T) {
		validUnixCompressionTypes := []any{
			"tgz", "zip",
		}

		for _, unixCompressionType := range validUnixCompressionTypes {
			_, err := NewUnixCompressionType(unixCompressionType)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'",
					unixCompressionType, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidUnixCompressionType", func(t *testing.T) {
		invalidUnixCompressionTypes := []any{
			"", "jpeg", "pdf",
		}

		for _, unixCompressionType := range invalidUnixCompressionTypes {
			_, err := NewUnixCompressionType(unixCompressionType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", unixCompressionType)
			}
		}
	})
}
