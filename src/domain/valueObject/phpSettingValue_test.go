package valueObject

import "testing"

func TestPhpSettingValue(t *testing.T) {
	t.Run("ValidPhpSettingValues", func(t *testing.T) {
		validValues := []interface{}{
			"on",
			"off",
			"ON",
			"OFF",
			"true",
			"false",
			"TRUE",
			"FALSE",
			true,
			false,
			0,
			1,
			2,
			"test",
			"dev",
			"prod",
		}

		for _, value := range validValues {
			_, err := NewPhpSettingValue(value)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", value, err)
			}
		}
	})

	t.Run("InvalidPhpSettingValues", func(t *testing.T) {
		invalidValues := []string{""}

		for _, value := range invalidValues {
			_, err := NewPhpSettingValue(value)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", value)
			}
		}
	})
}
