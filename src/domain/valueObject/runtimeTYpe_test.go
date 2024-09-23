package valueObject

import "testing"

func TestRuntimeType(t *testing.T) {
	t.Run("ValidRuntimeType", func(t *testing.T) {
		validDbTypes := []interface{}{
			"php-webserver", "php-ws", "php", "lsphp",
		}

		for _, dbType := range validDbTypes {
			_, err := NewRuntimeType(dbType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", dbType, err.Error())
			}
		}
	})

	t.Run("InvalidRuntimeType", func(t *testing.T) {
		invalidDbTypes := []interface{}{
			"tomcat", "python", "cosmosdb",
		}

		for _, dbType := range invalidDbTypes {
			_, err := NewRuntimeType(dbType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", dbType)
			}
		}
	})
}
