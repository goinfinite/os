package valueObject

import "testing"

func TestCronComment(t *testing.T) {
	t.Run("ValidCronComment", func(t *testing.T) {
		validCronComments := []interface{}{
			"Daily backup", "Database update at 3 AM",
			"Weekly report generated every Monday", "Temporary files cleanup",
		}

		for _, cronComment := range validCronComments {
			_, err := NewCronComment(cronComment)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", cronComment, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidCronComment", func(t *testing.T) {
		invalidCronComments := []interface{}{
			"A", "", nil,
		}

		for _, cronComment := range invalidCronComments {
			_, err := NewCronComment(cronComment)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", cronComment)
			}
		}
	})
}
