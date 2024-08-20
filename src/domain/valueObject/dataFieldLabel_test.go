package valueObject

import "testing"

func TestDataFieldLabel(t *testing.T) {
	t.Run("ValidDataFieldLabel", func(t *testing.T) {
		validDataFieldLabels := []interface{}{
			"Administrator Email", "Super administrator password",
			"Public directory to install", "Your own custom domain",
			"Container port binding", "Customem phone number",
		}

		for _, label := range validDataFieldLabels {
			_, err := NewDataFieldLabel(label)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", label, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldLabel", func(t *testing.T) {
		invalidDataFieldLabels := []interface{}{
			"", "./test", "-key", "anotherkey-",
		}

		for _, label := range invalidDataFieldLabels {
			_, err := NewDataFieldLabel(label)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", label)
			}
		}
	})
}
