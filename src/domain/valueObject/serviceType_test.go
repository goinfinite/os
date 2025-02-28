package valueObject

import (
	"slices"
	"testing"
)

func TestServiceType(t *testing.T) {
	t.Run("ValidServiceTypes", func(t *testing.T) {
		validServiceTypes := []string{
			"application", "mom", "monitoring", "logging", "security", "backup",
		}
		validServiceTypes = slices.Concat(validServiceTypes, ValidServiceTypes)

		for _, serviceType := range validServiceTypes {
			_, err := NewServiceType(serviceType)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", serviceType, err.Error(),
				)
			}
		}
	})
}
