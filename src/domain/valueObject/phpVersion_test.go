package valueObject

import "testing"

func TestPhpVersion(t *testing.T) {
	t.Run("ValidPhpVersions", func(t *testing.T) {
		validPhpVersions := []string{
			"1.0",
			"20",
		}
		for _, phpVersion := range validPhpVersions {
			_, err := NewPhpVersion(phpVersion)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", phpVersion, err)
			}
		}
	})

	t.Run("InvalidPhpVersions", func(t *testing.T) {
		invalidPhpVersions := []string{
			"1.0.0",
			"1.0.",
			"1..",
			"100",
		}
		for _, phpVersion := range invalidPhpVersions {
			_, err := NewPhpVersion(phpVersion)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", phpVersion)
			}
		}
	})
}
