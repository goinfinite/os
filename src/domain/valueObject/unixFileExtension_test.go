package valueObject

import "testing"

func TestUnixFileExtension(t *testing.T) {
	t.Run("ValidUnixFileExtension", func(t *testing.T) {
		validUnixFileExtensions := []interface{}{
			".png", "png", ".c", "c", ".ecelp4800", ".n-gage", ".application",
			".fe_launch", ".cdbcmsg",
		}

		for _, extension := range validUnixFileExtensions {
			_, err := NewUnixFileExtension(extension)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", extension, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidUnixFileExtension", func(t *testing.T) {
		invalidUnixFileExtensions := []interface{}{
			"", "file.php?blabla", "@<php52.sandbox.ntorga.com>.php", "../file.php",
			"hello10/info.php",
		}

		for _, extension := range invalidUnixFileExtensions {
			_, err := NewUnixFileExtension(extension)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", extension)
			}
		}
	})
}
