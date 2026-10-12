package valueObject

import "testing"

func TestMarketplaceItemSlug(t *testing.T) {
	t.Run("ValidMarketplaceItemSlug", func(t *testing.T) {
		validMarketplaceItemSlugs := []any{
			"drupal", "joomla", "lamp", "lemp", "laravel", "opencart", "oc",
			"wp", "wordpress",
		}

		for _, itemSlug := range validMarketplaceItemSlugs {
			_, err := NewMarketplaceItemSlug(itemSlug)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", itemSlug, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemSlug", func(t *testing.T) {
		invalidMarketplaceItemSlugs := []any{
			"", ".", "..", "/", "Slug with spaces", "<root>",
		}

		for _, itemSlug := range invalidMarketplaceItemSlugs {
			_, err := NewMarketplaceItemSlug(itemSlug)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", itemSlug)
			}
		}
	})
}
