package valueObject

import "testing"

func TestServiceName(t *testing.T) {
	t.Run("ValidServiceName", func(t *testing.T) {
		validServiceName := []interface{}{
			"php-webserver", "mariadb", "mysql", "postgresql", "python",
			"java", "nodejs", "python", "php", "redis-server",
		}

		for _, serviceName := range validServiceName {
			_, err := NewServiceName(serviceName)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", serviceName, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidServiceName", func(t *testing.T) {
		invalidServiceName := []interface{}{
			"nginx@", "my<>sql", "php#fpm", "node(js)",
		}

		for _, serviceName := range invalidServiceName {
			_, err := NewServiceName(serviceName)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", serviceName)
			}
		}
	})
}
