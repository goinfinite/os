package valueObject

import "testing"

func TestScheduledTaskId(t *testing.T) {
	t.Run("ValidScheduledTaskId", func(t *testing.T) {
		validScheduledTaskIds := []interface{}{
			0, 1, 445, "15987612309",
		}

		for _, taskId := range validScheduledTaskIds {
			_, err := NewScheduledTaskId(taskId)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", taskId, err.Error())
			}
		}
	})

	t.Run("InvalidScheduledTaskId", func(t *testing.T) {
		invalidScheduledTaskIds := []interface{}{
			-1, -455, "-15987612309",
		}

		for _, taskId := range invalidScheduledTaskIds {
			_, err := NewScheduledTaskId(taskId)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", taskId)
			}
		}
	})
}
