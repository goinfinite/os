package valueObject

import "testing"

func TestMarketplaceItemType(t *testing.T) {
	t.Run("ValidMarketplaceItemType", func(t *testing.T) {
		validMarketplaceItemTypes := []interface{}{
			"app", "framework", "stack",
		}

		for _, itemType := range validMarketplaceItemTypes {
			_, err := NewMarketplaceItemType(itemType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", itemType, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemType", func(t *testing.T) {
		invalidMarketplaceItemTypes := []interface{}{
			"", "service", "mobile", "ml-model", "repository",
		}

		for _, itemType := range invalidMarketplaceItemTypes {
			_, err := NewMarketplaceItemType(itemType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", itemType)
			}
		}
	})
}
