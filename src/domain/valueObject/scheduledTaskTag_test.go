package valueObject

import "testing"

func TestScheduledTaskTag(t *testing.T) {
	t.Run("ValidScheduledTaskTag", func(t *testing.T) {
		validScheduledTaskTag := []interface{}{
			"services", "marketplace", "ssl", "cron", "account",
		}

		for _, taskTag := range validScheduledTaskTag {
			_, err := NewScheduledTaskTag(taskTag)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", taskTag, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidScheduledTaskTag", func(t *testing.T) {
		invalidScheduledTaskTag := []interface{}{
			"", "123", "container!",
		}

		for _, taskTag := range invalidScheduledTaskTag {
			_, err := NewScheduledTaskTag(taskTag)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", taskTag)
			}
		}
	})
}
