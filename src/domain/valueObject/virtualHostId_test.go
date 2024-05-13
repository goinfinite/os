package valueObject

import (
	"testing"
)

func TestVirtualHostId(t *testing.T) {
	t.Run("ValidVirtualHostId", func(t *testing.T) {
		validVirtualHostIds := []interface{}{
			1,
			1000,
			65365,
			"12345",
		}

		for _, vhostId := range validVirtualHostIds {
			_, err := NewVirtualHostId(vhostId)
			if err != nil {
				t.Errorf("Expected no error for %v, got %s", vhostId, err.Error())
			}
		}
	})

	t.Run("InvalidVirtualHostId", func(t *testing.T) {
		invalidVirtualHostIds := []interface{}{
			-1,
			0,
			1000000000000000000,
			"-455",
		}

		for _, vhostId := range invalidVirtualHostIds {
			_, err := NewVirtualHostId(vhostId)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", vhostId)
			}
		}
	})
}
