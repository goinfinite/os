package valueObject

import "testing"

func TestUnixFilePermissions(t *testing.T) {
	t.Run("ValidUnixFilePermissions", func(t *testing.T) {
		validUnixFilePermissions := []string{
			"0",
			"5",
			"55",
			"644",
			"755",
			"1644",
			"4755",
		}
		for _, permissions := range validUnixFilePermissions {
			_, err := NewUnixFilePermissions(permissions)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", permissions, err)
			}
		}
	})

	t.Run("InvalidUnixFilePermissions", func(t *testing.T) {
		invalidUnixFilePermissions := []string{
			"-1",
			"",
			"00000",
			"aaaaa",
			"b",
		}
		for _, permissions := range invalidUnixFilePermissions {
			_, err := NewUnixFilePermissions(permissions)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", permissions)
			}
		}
	})

	t.Run("ValidUnixFilePermissionsFromInt", func(t *testing.T) {
		validUnixFileIntPermissions := []int{
			0,
			777,
			1777,
			7777,
		}
		for _, intPermissions := range validUnixFileIntPermissions {
			_, err := NewUnixFilePermissionsFromInt(intPermissions)
			if err != nil {
				t.Errorf("Expected no error for %d, got %v", intPermissions, err)
			}
		}
	})

	t.Run("InvalidUnixFilePermissionsFromInt", func(t *testing.T) {
		invalidUnixFileIntPermissions := []int{100000}
		for _, intPermissions := range invalidUnixFileIntPermissions {
			_, err := NewUnixFilePermissionsFromInt(intPermissions)
			if err == nil {
				t.Errorf("Expected error for %d, got nil", intPermissions)
			}
		}
	})
}
