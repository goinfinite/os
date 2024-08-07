package entity

import "testing"

func TestPhpSetting(t *testing.T) {
	t.Run("ValidPhpSetting", func(t *testing.T) {
		validPhpSettings := []string{
			"allow_url_fopen:on:firstOption", "file_uploads:OFF",
			"allow_url_include:ON:firstOption,secondOption",
			"date.timezone:America/Sao_Paulo:firstOption,secondOption,thidOption",
			"display_errors:off:firstOption,secondOption,thidOption,fourthOption",
		}

		for _, setting := range validPhpSettings {
			_, err := NewPhpSettingFromString(setting)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", setting, err.Error())
			}
		}
	})

	t.Run("InvalidPhpSetting", func(t *testing.T) {
		invalidPhpSettings := []string{
			"", "allow_url_fopen", "file_uploads:",
		}

		for _, setting := range invalidPhpSettings {
			_, err := NewPhpSettingFromString(setting)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", setting)
			}
		}
	})
}
