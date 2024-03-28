package valueObject

import "testing"

func TestMktplaceItemId(t *testing.T) {
	t.Run("ValidMktplaceItemId", func(t *testing.T) {
		validMktplaceItemIds := []interface{}{
			1,
			1000,
			65365,
			"12345",
		}

		for _, mii := range validMktplaceItemIds {
			_, err := NewMktplaceItemId(mii)
			if err != nil {
				t.Errorf("Expected no error for %v, got %s", mii, err.Error())
			}
		}
	})

	t.Run("InvalidMktplaceItemId", func(t *testing.T) {
		invalidMktplaceItemIds := []interface{}{
			-1,
			0,
			1000000000000000000,
			"-455",
		}

		for _, mii := range invalidMktplaceItemIds {
			_, err := NewMktplaceItemId(mii)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", mii)
			}
		}
	})
}
