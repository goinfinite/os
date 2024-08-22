package valueObject

import "testing"

func TestActivityRecordLevel(t *testing.T) {
	t.Run("ValidActivityRecordLevel", func(t *testing.T) {
		validActivityRecordLevels := []interface{}{
			"DEBUG", "INFO", "WARN", "ERROR", "SEC",
		}

		for _, level := range validActivityRecordLevels {
			_, err := NewActivityRecordLevel(level)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", level, err.Error())
			}
		}
	})

	t.Run("InvalidActivityRecordLevel", func(t *testing.T) {
		invalidActivityRecordLevels := []interface{}{
			"LOG", "MSG", "RECORD", "FYI",
		}

		for _, level := range invalidActivityRecordLevels {
			_, err := NewActivityRecordLevel(level)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", level)
			}
		}
	})
}
