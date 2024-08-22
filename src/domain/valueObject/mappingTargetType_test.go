package valueObject

import "testing"

func TestMappingTargetType(t *testing.T) {
	t.Run("ValidMappingTargetType", func(t *testing.T) {
		validMappingTargetTypes := []interface{}{
			"url", "service", "response-code", "inline-html", "static-files",
		}

		for _, targetType := range validMappingTargetTypes {
			_, err := NewMappingTargetType(targetType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", targetType, err.Error())
			}
		}
	})

	t.Run("InvalidMappingTargetType", func(t *testing.T) {
		invalidMappingTargetTypes := []interface{}{
			"response-header", "reverse-proxy", "template",
		}

		for _, targetType := range invalidMappingTargetTypes {
			_, err := NewMappingTargetType(targetType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", targetType)
			}
		}
	})
}
