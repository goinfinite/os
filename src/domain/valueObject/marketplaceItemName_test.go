package valueObject

import (
	"testing"
)

func TestMarketplaceItemName(t *testing.T) {
	t.Run("ValidMarketplaceItemName", func(t *testing.T) {
		validMarketplaceItemNames := []string{
			"wordpress",
			"WordPress",
			"opencart",
			"OpenCart",
			"Magento",
			"magento",
			"Joomla",
			"joomla",
			"Drupal",
			"drupal",
			"Supabase",
			"supabase",
			"Laravel",
			"laravel",
			"rabbitmq",
			"RabbitMQ",
		}
		for _, min := range validMarketplaceItemNames {
			_, err := NewMarketplaceItemName(min)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", min, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemName", func(t *testing.T) {
		invalidMarketplaceItemNames := []string{
			"",
			".",
			"..",
			"/",
			"name with space",
			"A very long name without any reason just for the test",
			"ççççççç",
			"<root>",
		}
		for _, min := range invalidMarketplaceItemNames {
			_, err := NewMarketplaceItemName(min)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", min)
			}
		}
	})
}
