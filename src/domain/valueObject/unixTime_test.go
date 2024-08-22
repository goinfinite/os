package valueObject

import (
	"testing"
	"time"
)

func TestUnixTime(t *testing.T) {
	t.Run("ValidUnixTime", func(t *testing.T) {
		validUnixTime := []interface{}{
			"0", int(0), int8(0), int16(0), int32(0), int64(0), uint(0), uint8(0),
			uint16(0), uint32(0), uint64(0), float32(0), float64(0),
		}

		for _, unixTime := range validUnixTime {
			_, err := NewUnixTime(unixTime)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", unixTime, err.Error(),
				)
			}
		}
	})

	t.Run("ValidUnixTimeNow", func(t *testing.T) {
		nowUnixTime := time.Now().Unix()
		validUnixTimeNow := NewUnixTimeNow().Int64()
		diff := nowUnixTime - validUnixTimeNow

		if diff != 0 {
			t.Errorf(
				"Expected no difference between '%d' and '%d', got '%d' as difference",
				nowUnixTime, validUnixTimeNow, diff,
			)
		}
	})
}
