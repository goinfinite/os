package presenterValueObject

import "testing"

func TestMarketplaceListType(t *testing.T) {
	t.Run("ValidMarketplaceListType", func(t *testing.T) {
		validMarketplaceListTypes := []interface{}{"installed", "catalog"}

		for _, listType := range validMarketplaceListTypes {
			_, err := NewMarketplaceListType(listType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", listType, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceListType", func(t *testing.T) {
		invalidMarketplaceListTypes := []interface{}{"cancelled", "uninstalled", "error"}

		for _, listType := range invalidMarketplaceListTypes {
			_, err := NewMarketplaceListType(listType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", listType)
			}
		}
	})
}
