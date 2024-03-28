package valueObject

import (
	"testing"
)

func TestMktplaceItemType(t *testing.T) {
	t.Run("ValidMktplaceItemType", func(t *testing.T) {
		validMktplaceItemTypes := []string{
			"app",
			"framework",
			"stack",
		}
		for _, mit := range validMktplaceItemTypes {
			_, err := NewMktplaceItemType(mit)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mit, err.Error())
			}
		}
	})

	t.Run("InvalidMktplaceItemType", func(t *testing.T) {
		invalidMktplaceItemTypes := []string{
			"",
			"service",
			"mobile",
			"ml-model",
			"repository",
		}
		for _, mit := range invalidMktplaceItemTypes {
			_, err := NewMktplaceItemType(mit)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mit)
			}
		}
	})
}
