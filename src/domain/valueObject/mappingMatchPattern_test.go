package valueObject

import "testing"

func TestMappingMatchPattern(t *testing.T) {
	t.Run("ValidMappingMatchPattern", func(t *testing.T) {
		validMappingMatchPatterns := []interface{}{
			"begins-with", "begins with", "contains", "equals", "ends-with",
			"ends with",
		}

		for _, matchPattern := range validMappingMatchPatterns {
			_, err := NewMappingMatchPattern(matchPattern)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", matchPattern, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidMappingMatchPattern", func(t *testing.T) {
		invalidMappingMatchPatterns := []interface{}{
			"", "bigger-then", "diff", "has-prefix",
		}

		for _, matchPattern := range invalidMappingMatchPatterns {
			_, err := NewMappingMatchPattern(matchPattern)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", matchPattern)
			}
		}
	})
}
