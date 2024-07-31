package valueObject

import "testing"

func TestCronId(t *testing.T) {
	t.Run("ValidCronId", func(t *testing.T) {
		validCronIds := []interface{}{
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

		for _, cronId := range validCronIds {
			_, err := NewCronId(cronId)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", cronId, err)
			}
		}
	})

	t.Run("InvalidCronId", func(t *testing.T) {
		invalidCronIds := []interface{}{
			"-1",
			int(-1),
			int8(-1),
			int16(-1),
			int32(-1),
			int64(-1),
			float32(-1),
			float64(-1),
		}

		for _, cronId := range invalidCronIds {
			_, err := NewCronId(cronId)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", cronId)
			}
		}
	})
}
