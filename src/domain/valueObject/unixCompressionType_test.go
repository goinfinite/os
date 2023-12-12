package valueObject

import "testing"

func TestUnixCompressionType(t *testing.T) {
	t.Run("ValidUnixCompressionType", func(t *testing.T) {
		validUnixCompressionTypes := []string{
			"gzip",
			"gzip",
		}
		for _, extension := range validUnixCompressionTypes {
			_, err := NewUnixCompressionType(extension)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", extension, err)
			}
		}
	})

	t.Run("InvalidUnixCompressionType", func(t *testing.T) {
		invalidUnixCompressionTypes := []string{
			"",
			"jpeg",
			"pdf",
		}
		for _, extension := range invalidUnixCompressionTypes {
			_, err := NewUnixCompressionType(extension)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", extension)
			}
		}
	})
}
