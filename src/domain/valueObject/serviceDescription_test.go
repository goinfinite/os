package valueObject

import "testing"

func TestServiceDescription(t *testing.T) {
	t.Run("ValidServiceDescription", func(t *testing.T) {
		validServiceDescription := []any{
			"php-webserver", "mariadb", "mysql", "postgresql", "python",
			"java", "nodejs", "python",
		}

		for _, serviceDescription := range validServiceDescription {
			_, err := NewServiceDescription(serviceDescription)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", serviceDescription, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidServiceDescription", func(t *testing.T) {
		invalidServiceDescription := []any{
			"", "a",
		}

		for _, serviceDescription := range invalidServiceDescription {
			_, err := NewServiceDescription(serviceDescription)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", serviceDescription)
			}
		}
	})
}
