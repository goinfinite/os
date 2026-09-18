package valueObject

import "testing"

func TestVirtualHostType(t *testing.T) {
	t.Run("ValidVirtualHostType", func(t *testing.T) {
		validVhostTypes := []any{
			"primary", "top-level", "subdomain", "wildcard", "alias",
		}

		for _, vhostType := range validVhostTypes {
			_, err := NewVirtualHostType(vhostType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", vhostType, err.Error())
			}
		}
	})

	t.Run("InvalidVirtualHostType", func(t *testing.T) {
		invalidVhostTypes := []any{
			"extradomain", "low-level", "secondary", "legacy",
		}

		for _, vhostType := range invalidVhostTypes {
			_, err := NewVirtualHostType(vhostType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", vhostType)
			}
		}
	})
}
