package valueObject

import (
	"testing"
)

func TestMappingTargetType(t *testing.T) {
	t.Run("ValidMappingTargetType", func(t *testing.T) {
		validMappingTargetTypes := []string{
			"url",
			"service",
			"response-code",
			"inline-html",
			"static-files",
		}

		for _, mtt := range validMappingTargetTypes {
			_, err := NewMappingTargetType(mtt)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mtt, err.Error())
			}
		}
	})

	t.Run("ValidMappingTargetType", func(t *testing.T) {
		invalidMappingTargetTypes := []string{
			"response-header",
			"reverse-proxy",
			"template",
		}

		for _, mtt := range invalidMappingTargetTypes {
			_, err := NewMappingTargetType(mtt)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mtt)
			}
		}
	})
}
