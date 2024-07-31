package valueObject

import "testing"

func TestMarketplaceItemId(t *testing.T) {
	t.Run("ValidMarketplaceItemId", func(t *testing.T) {
		validMarketplaceItemIds := []interface{}{
			0,
			1,
			3,
			1000,
			65365,
			"12345",
		}

		for _, itemId := range validMarketplaceItemIds {
			_, err := NewMarketplaceItemId(itemId)
			if err != nil {
				t.Errorf("(%v) ExpectedNoErrorButGot: %s", itemId, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemId", func(t *testing.T) {
		invalidMarketplaceItemIds := []interface{}{
			-1,
			9223372036854775807,
			"-455",
		}

		for _, itemId := range invalidMarketplaceItemIds {
			_, err := NewMarketplaceItemId(itemId)
			if err == nil {
				t.Errorf("(%v) ExpectedErrorButGotNil", itemId)
			}
		}
	})
}
