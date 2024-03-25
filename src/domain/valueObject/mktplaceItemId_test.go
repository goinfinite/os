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

		for _, groupId := range validMktplaceItemIds {
			_, err := NewMktplaceItemId(groupId)
			if err != nil {
				t.Errorf("Expected no error for %v, got %s", groupId, err.Error())
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

		for _, groupId := range invalidMktplaceItemIds {
			_, err := NewMktplaceItemId(groupId)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", groupId)
			}
		}
	})
}
