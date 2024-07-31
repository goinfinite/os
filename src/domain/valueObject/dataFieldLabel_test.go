package valueObject

import "testing"

func TestDataFieldLabel(t *testing.T) {
	t.Run("ValidDataFieldLabel", func(t *testing.T) {
		validDataFieldLabels := []string{
			"Administrator Email",
			"Super administrator password",
			"Public directory to install",
			"Your own custom domain",
			"Container port binding",
			"Customem phone number",
		}

		for _, dfl := range validDataFieldLabels {
			_, err := NewDataFieldLabel(dfl)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", dfl, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldLabel", func(t *testing.T) {
		invalidDataFieldLabels := []string{
			"",
			"./test",
			"-key",
			"anotherkey-",
		}

		for _, dfl := range invalidDataFieldLabels {
			_, err := NewDataFieldLabel(dfl)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dfl)
			}
		}
	})
}
