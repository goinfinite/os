package valueObject

import "testing"

func TestIpAddress(t *testing.T) {
	t.Run("ValidIpAddress", func(t *testing.T) {
		validIpAddresses := []interface{}{
			"192.168.1.1", "10.0.0.1", "172.16.0.1", "::1", "2001:db8::1",
		}

		for _, ipAddress := range validIpAddresses {
			_, err := NewIpAddress(ipAddress)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", ipAddress, err.Error())
			}
		}
	})

	t.Run("InvalidIpAddress", func(t *testing.T) {
		invalidIpAddresses := []interface{}{
			"192.168.1.256", "300.0.0.1", "123.456.78.90", "abcd::12345",
			"192.168.1.1.1",
		}

		for _, ipAddress := range invalidIpAddresses {
			_, err := NewIpAddress(ipAddress)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", ipAddress)
			}
		}
	})
}
