package valueObject

import (
	"testing"
)

func TestMarketplaceItemId(t *testing.T) {
	t.Run("ValidMarketplaceItemId", func(t *testing.T) {
		validMarketplaceItemIds := []interface{}{
			1,
			1000,
			65365,
			"12345",
		}

		for _, value := range validMarketplaceItemIds {
			_, err := NewMarketplaceItemId(value)
			if err != nil {
				t.Errorf("(%v) ExpectedNoErrorButGot: %s", value, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemId", func(t *testing.T) {
		invalidMarketplaceItemIds := []interface{}{
			-1,
			0,
			9223372036854775807,
			"-455",
		}

		for _, value := range invalidMarketplaceItemIds {
			_, err := NewMarketplaceItemId(value)
			if err == nil {
				t.Errorf("(%v) ExpectedErrorButGotNil", value)
			}
		}
	})
}
