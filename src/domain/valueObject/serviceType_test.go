package valueObject

import (
	"testing"
)

func TestServiceType(t *testing.T) {
	t.Run("ValidServiceTypes", func(t *testing.T) {
		for _, knownServiceType := range ValidServiceTypes {
			serviceType, err := NewServiceType(knownServiceType)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", serviceType, err.Error(),
				)
			}

			if serviceType.String() != knownServiceType {
				t.Errorf(
					"Expected '%v', got '%v'", knownServiceType, serviceType,
				)
			}
		}

		alsoValidServiceTypes := []string{
			"application", "backup", "logging", "mom", "monitoring", "other", "security",
		}
		for _, rawServiceType := range alsoValidServiceTypes {
			serviceType, err := NewServiceType(rawServiceType)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", serviceType, err.Error(),
				)
			}

			if serviceType.String() != "other" {
				t.Errorf(
					"Expected 'other', got '%v'", serviceType,
				)
			}
		}
	})
}
