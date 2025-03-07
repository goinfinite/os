package valueObject

import "testing"

func TestUnixFileOwnership(t *testing.T) {
	t.Run("ValidUnixFileOwnership", func(t *testing.T) {
		validUnixFileOwnerships := []interface{}{
			"dev:dev", "sudo:sudo", "root:root", "www-data:www-data", "www-data:root",
			"www-data:dev", "dev:www-data", "dev:root", "root:dev", "root:www-data",
		}

		for _, ownership := range validUnixFileOwnerships {
			_, err := NewUnixFileOwnership(ownership)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", ownership, err.Error())
			}
		}
	})

	t.Run("InvalidUnixFileOwnership", func(t *testing.T) {
		invalidUnixFileOwnerships := []interface{}{
			"", 1000, true, ":dev", "dev:", "dev:dev:dev", "dev/dev",
		}

		for _, ownership := range invalidUnixFileOwnerships {
			_, err := NewUnixFileOwnership(ownership)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", ownership)
			}
		}
	})
}
