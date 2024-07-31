package valueObject

import "testing"

func TestVirtualHostType(t *testing.T) {
	t.Run("ValidVirtualHostType", func(t *testing.T) {
		validVirtualHostTypes := []string{
			"primary", "top-level", "subdomain", "wildcard", "alias",
		}

		for _, vhostType := range validVirtualHostTypes {
			_, err := NewVirtualHostType(vhostType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", vhostType, err.Error())
			}
		}
	})

	t.Run("InvalidVirtualHostType", func(t *testing.T) {
		invalidVirtualHostTypes := []string{
			"secondary", "low-level", "domain", "target",
		}

		for _, vhostType := range invalidVirtualHostTypes {
			_, err := NewVirtualHostType(vhostType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", vhostType)
			}
		}
	})
}
