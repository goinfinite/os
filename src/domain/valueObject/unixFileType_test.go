package valueObject

import "testing"

func TestUnixFileType(t *testing.T) {
	t.Run("ValidUnixFileType", func(t *testing.T) {
		validUnixFileTypes := []string{
			"directory",
			"file",
		}
		for _, fileType := range validUnixFileTypes {
			_, err := NewUnixFileType(fileType)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", fileType, err)
			}
		}
	})

	t.Run("InvalidUnixFileType", func(t *testing.T) {
		invalidUnixFileTypes := []string{
			"",
			"jpg",
			"abcd",
			"aaaa222333",
		}
		for _, fileType := range invalidUnixFileTypes {
			_, err := NewUnixFileType(fileType)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", fileType)
			}
		}
	})
}
