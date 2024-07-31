package valueObject

import "testing"

func TestAccountId(t *testing.T) {
	t.Run("ValidAccountId", func(t *testing.T) {
		validAccountIds := []interface{}{
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

		for _, accountId := range validAccountIds {
			_, err := NewAccountId(accountId)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", accountId, err)
			}
		}
	})

	t.Run("InvalidAccountId", func(t *testing.T) {
		invalidAccountIds := []interface{}{
			"-1",
			int(-1),
			int8(-1),
			int16(-1),
			int32(-1),
			int64(-1),
			float32(-1),
			float64(-1),
		}

		for _, accountId := range invalidAccountIds {
			_, err := NewAccountId(accountId)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", accountId)
			}
		}
	})
}
