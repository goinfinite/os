package valueObject

import "testing"

func TestMarketplaceItemManifestVersion(t *testing.T) {
	t.Run("ValidMarketplaceItemManifestVersion", func(t *testing.T) {
		validMarketplaceItemManifestVersions := []interface{}{
			"v1",
		}

		for _, itemManifestVersion := range validMarketplaceItemManifestVersions {
			_, err := NewMarketplaceItemManifestVersion(itemManifestVersion)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'",
					itemManifestVersion, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidMarketplaceItemManifestVersion", func(t *testing.T) {
		invalidMarketplaceItemManifestVersions := []interface{}{
			"v0", 0, false, 1.00,
		}

		for _, itemManifestVersion := range invalidMarketplaceItemManifestVersions {
			_, err := NewMarketplaceItemManifestVersion(itemManifestVersion)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", itemManifestVersion)
			}
		}
	})
}
