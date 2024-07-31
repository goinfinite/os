package valueObject

import "testing"

func TestMarketplaceItemType(t *testing.T) {
	t.Run("ValidMarketplaceItemType", func(t *testing.T) {
		validMarketplaceItemTypes := []interface{}{
			"app",
			"framework",
			"stack",
		}
		for _, mit := range validMarketplaceItemTypes {
			_, err := NewMarketplaceItemType(mit)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mit, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemType", func(t *testing.T) {
		invalidMarketplaceItemTypes := []interface{}{
			"",
			"service",
			"mobile",
			"ml-model",
			"repository",
		}
		for _, mit := range invalidMarketplaceItemTypes {
			_, err := NewMarketplaceItemType(mit)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mit)
			}
		}
	})
}
