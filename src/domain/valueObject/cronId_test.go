package valueObject

import (
	"testing"
)

func TestNewCronId(t *testing.T) {
	t.Run("ValidId", func(t *testing.T) {
		validIds := []interface{}{
			"0",
			int(0),
			int8(0),
			int16(0),
			int32(0),
			int64(0),
			uint(0),
			uint8(0),
			uint16(0),
			uint32(0),
			uint64(0),
			float32(0),
			float64(0),
		}

		for _, schedule := range validIds {
			_, err := NewCronId(schedule)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", schedule, err)
			}
		}
	})

	t.Run("InvalidId", func(t *testing.T) {
		invalidIds := []interface{}{
			"-1",
			int(-1),
			int8(-1),
			int16(-1),
			int32(-1),
			int64(-1),
			float32(-1),
			float64(-1),
		}

		for _, schedule := range invalidIds {
			_, err := NewCronId(schedule)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", schedule)
			}
		}
	})
}
