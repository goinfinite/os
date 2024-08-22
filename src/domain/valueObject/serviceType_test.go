package valueObject

import "testing"

func TestServiceType(t *testing.T) {
	t.Run("ValidServiceTypes", func(t *testing.T) {
		validServiceTypes := []interface{}{
			"application", "runtime", "database", "webserver", "mom", "monitoring",
			"logging", "security", "backup", "system", "other",
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

	t.Run("InvalidServiceTypees", func(t *testing.T) {
		invalidServiceTypes := []interface{}{
			"framework", "erp", "crm", "ai", "cloud",
		}

		for _, serviceType := range invalidServiceTypes {
			_, err := NewServiceType(serviceType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", serviceType)
			}
		}
	})
}
