package valueObject

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestMappingMatchPattern(t *testing.T) {
	t.Run("ValidMappingMatchPattern", func(t *testing.T) {
		validMappingMatchPatterns := []string{
			"begins-with",
			"contains",
			"equals",
			"ends-with",
		}

		for _, mmp := range validMappingMatchPatterns {
			_, err := NewMappingMatchPattern(mmp)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mmp, err.Error())
			}
		}
	})

	t.Run("InvalidMappingMatchPattern", func(t *testing.T) {
		invalidLength := 70
		invalidMappingMatchPatterns := []string{
			"",
			"bigger-then",
			"diff",
			"has-prefix",
			testHelpers.GenerateString(invalidLength),
		}

		for _, mmp := range invalidMappingMatchPatterns {
			_, err := NewMappingMatchPattern(mmp)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mmp)
			}
		}
	})
}
