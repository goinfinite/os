package valueObject

import "testing"

func TestMktplaceItemName(t *testing.T) {
	t.Run("ValidMktplaceItemName", func(t *testing.T) {
		validMktplaceItemNames := []string{
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
		for _, mktplaceItemName := range validMktplaceItemNames {
			_, err := NewMktplaceItemName(mktplaceItemName)
			if err != nil {
				t.Errorf(
					"Expected no error for %s, got %s",
					mktplaceItemName,
					err.Error(),
				)
			}
		}
	})

	t.Run("ValidMktplaceItemName", func(t *testing.T) {
		invalidMktplaceItemNames := []string{
			"",
			".",
			"..",
			"/",
			"name with space",
			"A very long name without any reason just for the test",
			"ççççççç",
			"<root>",
		}
		for _, mktplaceItemName := range invalidMktplaceItemNames {
			_, err := NewMktplaceItemName(mktplaceItemName)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mktplaceItemName)
			}
		}
	})
}
