package valueObject

import (
	"testing"
)

func TestSecureAccessPublicKeyId(t *testing.T) {
	t.Run("ValidSecureAccessPublicKeyId", func(t *testing.T) {
		rawValidSecureAccessPublicKeyId := []interface{}{
			"0", int(0), int8(0), int16(0), int32(0), int64(0), uint(0), uint8(0),
			uint16(0), uint32(0), uint64(0), float32(0), float64(0),
		}

		for _, rawKeyId := range rawValidSecureAccessPublicKeyId {
			_, err := NewSecureAccessPublicKeyId(rawKeyId)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", rawKeyId, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSecureAccessPublicKeyId", func(t *testing.T) {
		rawInvalidSecureAccessPublicKeyId := []interface{}{
			"-1", int(-1), int8(-1), int16(-1), int32(-1), int64(-1), float32(-1),
			float64(-1),
		}

		for _, rawKeyId := range rawInvalidSecureAccessPublicKeyId {
			_, err := NewSecureAccessPublicKeyId(rawKeyId)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", rawKeyId)
			}
		}
	})
}
