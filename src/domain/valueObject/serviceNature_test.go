package valueObject

import "testing"

func TestServiceNature(t *testing.T) {
	t.Run("ValidServiceNature", func(t *testing.T) {
		for _, serviceNature := range ValidServiceNatures {
			_, err := NewServiceNature(serviceNature)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'",
					serviceNature, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidServiceNature", func(t *testing.T) {
		invalidServiceNature := []interface{}{
			"installable", "executable", "downloadable",
		}

		for _, serviceNature := range invalidServiceNature {
			_, err := NewServiceNature(serviceNature)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", serviceNature)
			}
		}
	})
}
