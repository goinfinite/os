package valueObject

import (
	"testing"
)

func TestNewByte(t *testing.T) {
	t.Run("ValidByte", func(t *testing.T) {
		validBytes := []interface{}{
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

		for _, byteValue := range validBytes {
			_, err := NewByte(byteValue)
			if err != nil {
				t.Errorf("Expected no error for %v, got %v", byteValue, err)
			}
		}
	})

	t.Run("InvalidByte", func(t *testing.T) {
		invalidBytes := []interface{}{
			"-1",
			int(-1),
			int8(-1),
			int16(-1),
			int32(-1),
			int64(-1),
			float32(-1),
			float64(-1),
		}

		for _, byteValue := range invalidBytes {
			_, err := NewByte(byteValue)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", byteValue)
			}
		}
	})
}
