package valueObject

import "testing"

func TestServiceManifestVersion(t *testing.T) {
	t.Run("ValidServiceManifestVersion", func(t *testing.T) {
		validServiceManifestVersions := []any{
			"v1",
		}

		for _, manifestVersion := range validServiceManifestVersions {
			_, err := NewServiceManifestVersion(manifestVersion)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", manifestVersion, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidServiceManifestVersion", func(t *testing.T) {
		invalidServiceManifestVersions := []any{
			"v0", 0, false, 1.00,
		}

		for _, manifestVersion := range invalidServiceManifestVersions {
			_, err := NewServiceManifestVersion(manifestVersion)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", manifestVersion)
			}
		}
	})
}
