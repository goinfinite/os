package valueObject

import "testing"

func TestScheduledTaskStatus(t *testing.T) {
	t.Run("ValidScheduledTaskStatus", func(t *testing.T) {
		validScheduledTaskStatus := []any{
			"pending", "running", "completed", "failed", "cancelled", "timeout",
		}

		for _, taskStatus := range validScheduledTaskStatus {
			_, err := NewScheduledTaskStatus(taskStatus)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", taskStatus, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidScheduledTaskStatus", func(t *testing.T) {
		invalidScheduledTaskStatus := []any{
			"started", "success", "error",
		}

		for _, taskStatus := range invalidScheduledTaskStatus {
			_, err := NewScheduledTaskStatus(taskStatus)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", taskStatus)
			}
		}
	})
}
