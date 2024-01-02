package valueObject

import "testing"

func TestUnixCompressionType(t *testing.T) {
	t.Run("ValidUnixCompressionType", func(t *testing.T) {
		validUnixCompressionTypes := []string{
			"gzip",
			"zip",
		}
		for _, unixCompressionType := range validUnixCompressionTypes {
			_, err := NewUnixCompressionType(unixCompressionType)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", unixCompressionType, err)
			}
		}
	})

	t.Run("InvalidUnixCompressionType", func(t *testing.T) {
		invalidUnixCompressionTypes := []string{
			"",
			"jpeg",
			"pdf",
		}
		for _, unixCompressionType := range invalidUnixCompressionTypes {
			_, err := NewUnixCompressionType(unixCompressionType)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", unixCompressionType)
			}
		}
	})
}
