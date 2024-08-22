package valueObject

import "testing"

func TestScheduledTaskOutput(t *testing.T) {
	t.Run("ValidScheduledTaskOutput", func(t *testing.T) {
		validScheduledTaskOutputs := []interface{}{
			"validOutput", "everythingIsUpToDate", "ExecutedWithoutAnyError",
		}

		for _, taskOutput := range validScheduledTaskOutputs {
			_, err := NewScheduledTaskOutput(taskOutput)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", taskOutput, err.Error(),
				)
			}
		}
	})
}
