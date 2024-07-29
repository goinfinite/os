package valueObject

import (
	"testing"
)

func TestCronSchedule(t *testing.T) {
	t.Run("ValidSchedule", func(t *testing.T) {
		validSchedules := []string{
			"@annually",
			"@yearly",
			"@monthly",
			"@weekly",
			"@daily",
			"* * * * *",
			"*/5 * * * *",
			"0 0 * * *",
			"0 0 1 * *",
			"0,15,45 0 * * 0",
			"0 0 * * 0",
			"0 0 * * 1-5",
			"*/5 * * * 1-5",
			"@every 1h30m",
			"@every 5s",
		}

		for _, schedule := range validSchedules {
			_, err := NewCronSchedule(schedule)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", schedule, err)
			}
		}
	})

	t.Run("InvalidSchedule", func(t *testing.T) {
		invalidSchedules := []string{
			"*/5 * * * * *",
			"0 0 * *",
			"0 0 * * 0 0",
			"0 0 * * 0 0 0",
			"0 0 * * 0 0 0 0",
			"@every 1h30",
			"@blabla",
			"<script>alert('xss')</script>",
		}

		for _, schedule := range invalidSchedules {
			_, err := NewCronSchedule(schedule)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", schedule)
			}
		}
	})
}
