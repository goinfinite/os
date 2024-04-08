package valueObject

import (
	"testing"
)

func TestMarketplaceItemId(t *testing.T) {
	t.Run("ValidMarketplaceItemId", func(t *testing.T) {
		validMarketplaceItemIds := []interface{}{
			1,
			1000,
			65365,
			"12345",
		}

		for _, mii := range validMarketplaceItemIds {
			_, err := NewMarketplaceItemId(mii)
			if err != nil {
				t.Errorf("Expected no error for %v, got %s", mii, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemId", func(t *testing.T) {
		invalidMarketplaceItemIds := []interface{}{
			-1,
			0,
			1000000000000000000,
			"-455",
		}

		for _, mii := range invalidMarketplaceItemIds {
			_, err := NewMarketplaceItemId(mii)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", mii)
			}
		}
	})
}
