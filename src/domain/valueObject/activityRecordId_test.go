package valueObject

import "testing"

func TestActivityRecordId(t *testing.T) {
	t.Run("ValidActivityRecordId", func(t *testing.T) {
		validActivityRecordIds := []interface{}{
			"0", int(0), int8(0), int16(0), int32(0), int64(0), uint(0), uint8(0),
			uint16(0), uint32(0), uint64(0), float32(0), float64(0),
		}

		for _, id := range validActivityRecordIds {
			_, err := NewActivityRecordId(id)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", id, err.Error())
			}
		}
	})

	t.Run("InvalidActivityRecordId", func(t *testing.T) {
		invalidActivityRecordIds := []interface{}{
			"-1", int(-1), int8(-1), int16(-1), int32(-1), int64(-1), float32(-1),
			float64(-1),
		}

		for _, id := range invalidActivityRecordIds {
			_, err := NewActivityRecordId(id)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", id)
			}
		}
	})
}
