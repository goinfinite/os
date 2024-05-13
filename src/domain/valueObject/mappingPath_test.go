package valueObject

import (
	"testing"
)

func TestNewMappingPath(t *testing.T) {
	t.Run("ValidMappingPath", func(t *testing.T) {
		validMappingPaths := []string{
			"",
			"/",
			"/img/",
			"/index.html",
			".(png|gif|ico|jpg|jpeg)",
			"/(media|images|cache|tmp|logs)/.*.(php|jsp|pl|py|asp|cgi|sh)$",
			"something",
			"@opencart",
			"/(uploads|files|wp-content|wp-includes|akismet)/.*.php",
			"\\.php(/|$)",
		}

		for _, path := range validMappingPaths {
			_, err := NewMappingPath(path)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", path, err.Error())
			}
		}
	})

	t.Run("InvalidMappingPath", func(t *testing.T) {
		invalidMappingPaths := []string{
			"UNION SELECT * FROM USERS",
			"/path\n/path",
			"?param=value",
			"https://www.google.com",
			"/path/'; DROP TABLE users; --",
		}

		for _, path := range invalidMappingPaths {
			_, err := NewMappingPath(path)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", path)
			}
		}
	})
}
