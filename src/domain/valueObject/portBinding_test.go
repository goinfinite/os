package valueObject

import "testing"

func TestPortBinding(t *testing.T) {
	t.Run("ValidPortBindings", func(t *testing.T) {
		validPortBindings := []interface{}{
			80, "80/http", "443/https", "3306/tcp", "8000",
		}

		for _, portBinding := range validPortBindings {
			_, err := NewPortBinding(portBinding)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", portBinding, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidPortBindings", func(t *testing.T) {
		invalidPortBindings := []interface{}{
			"", "8000/",
		}

		for _, portBinding := range invalidPortBindings {
			_, err := NewPortBinding(portBinding)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", portBinding)
			}
		}
	})
}
