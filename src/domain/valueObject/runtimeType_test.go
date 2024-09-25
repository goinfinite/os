package valueObject

import "testing"

func TestRuntimeType(t *testing.T) {
	t.Run("ValidRuntimeType", func(t *testing.T) {
		validRuntimeTypes := []interface{}{
			"php-webserver", "php-ws", "php", "lsphp",
		}

		for _, runtimeType := range validRuntimeTypes {
			_, err := NewRuntimeType(runtimeType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", runtimeType, err.Error())
			}
		}
	})

	t.Run("InvalidRuntimeType", func(t *testing.T) {
		invalidRuntimeTypes := []interface{}{
			"tomcat", "python", "cosmosdb",
		}

		for _, runtimeType := range invalidRuntimeTypes {
			_, err := NewRuntimeType(runtimeType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", runtimeType)
			}
		}
	})
}
