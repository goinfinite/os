package valueObject

import "testing"

func TestServiceNameWithVersion(t *testing.T) {
	t.Run("ValidServiceNameWithVersion (from string)", func(t *testing.T) {
		validServiceNameWithVersion := []interface{}{
			"php-webserver:8.0", "mariadb:latest", "mysql:lts", "postgresql:alpha",
			"python:alpha", "java:beta", "nodejs", "redis-server:1.0.0", "nginx:1",
		}

		for _, serviceNameWithVersion := range validServiceNameWithVersion {
			_, err := NewServiceNameWithVersionFromString(serviceNameWithVersion)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'",
					serviceNameWithVersion, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidServiceNameWithVersion (from string)", func(t *testing.T) {
		invalidServiceNameWithVersion := []interface{}{
			"", "my<>sql", "nodejs:1.0<0", "mysql:",
		}

		for _, serviceNameWithVersion := range invalidServiceNameWithVersion {
			_, err := NewServiceNameWithVersionFromString(serviceNameWithVersion)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", serviceNameWithVersion)
			}
		}
	})
}
