package valueObject

import "testing"

func TestScheduledTaskName(t *testing.T) {
	t.Run("ValidScheduledTaskName", func(t *testing.T) {
		validScheduledTaskNames := []interface{}{
			"installWordpress", "CreateCronTaskWhenOsInitialize",
			"testAllComponentsBeforeStart", "CheckIfPort443UsingSelfSignedSsl",
		}

		for _, taskName := range validScheduledTaskNames {
			_, err := NewScheduledTaskName(taskName)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", taskName, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidScheduledTaskName", func(t *testing.T) {
		invalidScheduledTaskNames := []interface{}{
			"", "1failedRequest", "ValidateUserInput!",
		}

		for _, taskName := range invalidScheduledTaskNames {
			_, err := NewScheduledTaskName(taskName)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", taskName)
			}
		}
	})
}
