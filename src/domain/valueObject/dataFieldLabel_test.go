package valueObject

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestDataFieldLabel(t *testing.T) {
	t.Run("ValidDataFieldLabel", func(t *testing.T) {
		validDataFieldLabels := []string{
			"username",
			"user-email",
			"Service-Name_With_Port80",
		}

		for _, dfl := range validDataFieldLabels {
			_, err := NewDataFieldLabel(dfl)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", dfl, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldLabel", func(t *testing.T) {
		invalidLength := 70
		invalidDataFieldLabels := []string{
			"",
			"./test",
			"-key",
			"anotherkey-",
			testHelpers.GenerateString(invalidLength),
		}

		for _, dfl := range invalidDataFieldLabels {
			_, err := NewDataFieldLabel(dfl)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dfl)
			}
		}
	})
}
