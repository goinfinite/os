package valueObject

import "testing"

func TestNetworkPort(t *testing.T) {
	t.Run("ValidNetworkPort", func(t *testing.T) {
		validNetworkPorts := []interface{}{
			"8080",
			int(443),
			int8(80),
			int16(8000),
			int32(8080),
			int64(8443),
			uint(443),
			uint8(80),
			uint16(8000),
			uint32(8080),
			uint64(8443),
			float32(8080),
			float64(8443),
		}

		for _, networkPort := range validNetworkPorts {
			_, err := NewNetworkPort(networkPort)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", networkPort, err)
			}
		}
	})

	t.Run("InvalidNetworkPort", func(t *testing.T) {
		invalidNetworkPorts := []interface{}{
			"-1",
			int(-1),
			int8(-1),
			int16(-1),
			int32(-1),
			int64(-1),
			float32(-1),
			float64(-1),
		}

		for _, networkPort := range invalidNetworkPorts {
			_, err := NewNetworkPort(networkPort)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", networkPort)
			}
		}
	})
}
