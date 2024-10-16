package valueObject

import "testing"

func TestNetworkProtocol(t *testing.T) {
	t.Run("ValidNetworkProtocol", func(t *testing.T) {
		validNetworkProtocols := []interface{}{
			"http", "https", "ws", "wss", "grpc", "grpcs", "tcp", "udp",
		}

		for _, networkProtocol := range validNetworkProtocols {
			_, err := NewNetworkProtocol(networkProtocol)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", networkProtocol, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidNetworkProtocol", func(t *testing.T) {
		invalidNetworkProtocols := []interface{}{
			"", "ftp", "dhcp", "smtp",
		}

		for _, networkProtocol := range invalidNetworkProtocols {
			_, err := NewNetworkProtocol(networkProtocol)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", networkProtocol)
			}
		}
	})
}
