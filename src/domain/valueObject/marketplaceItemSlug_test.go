package valueObject

import "testing"

func TestMarketplaceItemSlug(t *testing.T) {
	t.Run("ValidMarketplaceItemSlug", func(t *testing.T) {
		validMarketplaceItemSlugs := []interface{}{
			"drupal",
			"joomla",
			"lamp",
			"lemp",
			"laravel",
			"opencart",
			"oc",
			"wp",
			"wordpress",
		}
		for _, itemSlug := range validMarketplaceItemSlugs {
			_, err := NewMarketplaceItemSlug(itemSlug)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", itemSlug, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemSlug", func(t *testing.T) {
		invalidMarketplaceItemSlugs := []interface{}{
			"",
			".",
			"..",
			"/",
			"Slug with spaces",
			"<root>",
		}
		for _, itemSlug := range invalidMarketplaceItemSlugs {
			_, err := NewMarketplaceItemSlug(itemSlug)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", itemSlug)
			}
		}
	})
}
