package valueObject

import "testing"

func TestPhpVersion(t *testing.T) {
	t.Run("ValidPhpVersions", func(t *testing.T) {
		validPhpVersions := []interface{}{
			"1.0", "20",
		}

		for _, phpVersion := range validPhpVersions {
			_, err := NewPhpVersion(phpVersion)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", phpVersion, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidPhpVersions", func(t *testing.T) {
		invalidPhpVersions := []interface{}{
			"1.0.0", "1.0.", "1..", "100",
		}

		for _, phpVersion := range invalidPhpVersions {
			_, err := NewPhpVersion(phpVersion)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", phpVersion)
			}
		}
	})
}
