package valueObject

import "testing"

func TestMarketplaceItemName(t *testing.T) {
	t.Run("ValidMarketplaceItemName", func(t *testing.T) {
		validMarketplaceItemNames := []interface{}{
			"wordpress", "WordPress", "opencart", "OpenCart", "Magento", "magento",
			"Joomla", "joomla", "Drupal", "drupal", "Supabase", "supabase",
			"Laravel", "laravel", "rabbitmq", "RabbitMQ",
		}

		for _, name := range validMarketplaceItemNames {
			_, err := NewMarketplaceItemName(name)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", name, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemName", func(t *testing.T) {
		invalidMarketplaceItemNames := []interface{}{
			"", ".", "..", "/", "A very long name without any reason just for the test",
			"<root>",
		}

		for _, name := range invalidMarketplaceItemNames {
			_, err := NewMarketplaceItemName(name)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", name)
			}
		}
	})
}
