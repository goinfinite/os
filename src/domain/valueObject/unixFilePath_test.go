package valueObject

import "testing"

func TestUnixFilePath(t *testing.T) {
	t.Run("ValidUnixFilePath", func(t *testing.T) {
		validUnixFilePaths := []string{
			"/speedia/ssl_crt.pem",
			"/speedia/ssl_key.pem",
			"/usr/local/test.sh",
			"/speedia/ssl_crt",
		}
		for _, name := range validUnixFilePaths {
			_, err := NewUnixFilePath(name)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", name, err)
			}
		}
	})
}
