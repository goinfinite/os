package valueObject

import "testing"

func TestServiceEnv(t *testing.T) {
	t.Run("ValidServiceEnv", func(t *testing.T) {
		validServiceEnv := []interface{}{
			"NODE_ENV=development", "LOG_LEVEL=debug",
			"RUN_IN_BACKGROUND_MODE=true",
		}

		for _, serviceEnv := range validServiceEnv {
			_, err := NewServiceEnv(serviceEnv)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", serviceEnv, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidServiceEnv", func(t *testing.T) {
		invalidServiceEnv := []interface{}{
			"", "=development", "LOG_LEVEL=", "=", "NODE_ENV", true,
		}

		for _, serviceEnv := range invalidServiceEnv {
			_, err := NewServiceEnv(serviceEnv)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", serviceEnv)
			}
		}
	})
}
