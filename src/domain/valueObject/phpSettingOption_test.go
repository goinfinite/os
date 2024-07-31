package valueObject

import "testing"

func TestPhpSettingOption(t *testing.T) {
	t.Run("ValidPhpSettingOptions", func(t *testing.T) {
		validSettingOptions := []interface{}{
			"allow_url_fopen", "allow_url_include", "date.timezone", "display_errors",
			"error_log",
		}

		for _, name := range validSettingOptions {
			_, err := NewPhpSettingOption(name)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", name, err.Error())
			}
		}
	})

	t.Run("InvalidPhpSettingOptions", func(t *testing.T) {
		invalidSettingOptions := []interface{}{
			"",
		}

		for _, name := range invalidSettingOptions {
			_, err := NewPhpSettingOption(name)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", name)
			}
		}
	})
}
