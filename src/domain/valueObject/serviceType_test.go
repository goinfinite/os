package valueObject

import (
	"testing"
)

func TestServiceType(t *testing.T) {
	t.Run("ValidServiceTypes", func(t *testing.T) {
		validServiceTypes := []string{
			"application", "backup", "database", "logging", "mom", "monitoring",
			"other", "runtime", "security", "webserver",
		}

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
