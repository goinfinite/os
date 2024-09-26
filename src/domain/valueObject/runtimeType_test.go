package valueObject

import "testing"

func TestRuntimeType(t *testing.T) {
	t.Run("ValidRuntimeType", func(t *testing.T) {
		validRuntimeTypes := []interface{}{
			"php",
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
			"jre", "nodejs",
		}

		for _, runtimeType := range invalidRuntimeTypes {
			_, err := NewRuntimeType(runtimeType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", runtimeType)
			}
		}
	})
}
