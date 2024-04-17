package valueObject

import (
	"testing"
)

func TestMarketplaceInstalledItemId(t *testing.T) {
	t.Run("ValidMarketplaceInstalledItemId", func(t *testing.T) {
		validMarketplaceInstalledItemIds := []interface{}{
			1,
			1000,
			65365,
			"12345",
		}

		for _, miii := range validMarketplaceInstalledItemIds {
			_, err := NewMarketplaceInstalledItemId(miii)
			if err != nil {
				t.Errorf("Expected no error for %v, got %s", miii, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceInstalledItemId", func(t *testing.T) {
		invalidMarketplaceInstalledItemIds := []interface{}{
			-1,
			0,
			1000000000000000000,
			"-455",
		}

		for _, miii := range invalidMarketplaceInstalledItemIds {
			_, err := NewMarketplaceInstalledItemId(miii)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", miii)
			}
		}
	})
}
