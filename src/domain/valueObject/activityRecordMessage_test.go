package valueObject

import "testing"

func TestActivityRecordMessage(t *testing.T) {
	t.Run("ValidActivityRecordMessage", func(t *testing.T) {
		validActivityRecordMessages := []interface{}{
			"Something went wrong with respective scheduled task",
			"Unable to install marketplace item", "Error with install PHP",
		}

		for _, message := range validActivityRecordMessages {
			_, err := NewActivityRecordMessage(message)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", message, err.Error())
			}
		}
	})
}
