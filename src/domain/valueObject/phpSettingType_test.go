package valueObject

import "testing"

func TestPhpSettingType(t *testing.T) {
	t.Run("ValidPhpSettingType", func(t *testing.T) {
		validPhpSettingTypes := []interface{}{
			"select", "text",
		}

		for _, phpSettingType := range validPhpSettingTypes {
			_, err := NewPhpSettingType(phpSettingType)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", phpSettingType, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidPhpSettingType", func(t *testing.T) {
		invalidPhpSettingTypes := []interface{}{
			"button", "checkbox", "datetime-local", "date", "file", "hidden", "month",
			"reset", "submit", "week",
		}

		for _, phpSettingType := range invalidPhpSettingTypes {
			_, err := NewPhpSettingType(phpSettingType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", phpSettingType)
			}
		}
	})
}
