package valueObject

import (
	"testing"
)

func TestMarketplaceCatalogItemId(t *testing.T) {
	t.Run("ValidMarketplaceCatalogItemId", func(t *testing.T) {
		validMarketplaceCatalogItemIds := []interface{}{
			1,
			1000,
			65365,
			"12345",
		}

		for _, mcii := range validMarketplaceCatalogItemIds {
			_, err := NewMarketplaceCatalogItemId(mcii)
			if err != nil {
				t.Errorf("Expected no error for %v, got %s", mcii, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceCatalogItemId", func(t *testing.T) {
		invalidMarketplaceCatalogItemIds := []interface{}{
			-1,
			0,
			1000000000000000000,
			"-455",
		}

		for _, mcii := range invalidMarketplaceCatalogItemIds {
			_, err := NewMarketplaceCatalogItemId(mcii)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", mcii)
			}
		}
	})
}
