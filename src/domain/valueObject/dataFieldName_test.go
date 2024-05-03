package valueObject

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestDataFieldKey(t *testing.T) {
	t.Run("ValidDataFieldKey", func(t *testing.T) {
		validDataFieldKeys := []string{
			"username",
			"user-email",
			"Service-Name_With_Port80",
		}

		for _, dfk := range validDataFieldKeys {
			_, err := NewDataFieldName(dfk)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", dfk, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldKey", func(t *testing.T) {
		invalidLength := 70
		invalidDataFieldKeys := []string{
			"",
			"./test",
			"-key",
			"anotherkey-",
			testHelpers.GenerateString(invalidLength),
		}

		for _, dfk := range invalidDataFieldKeys {
			_, err := NewDataFieldName(dfk)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dfk)
			}
		}
	})
}
