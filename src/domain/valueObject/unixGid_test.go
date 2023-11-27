package valueObject

import "testing"

func TestUnixGid(t *testing.T) {
	t.Run("ValidUnixGid", func(t *testing.T) {
		validUnixGids := []int{0, 1000, 65365}
		for _, unixGid := range validUnixGids {
			_, err := NewUnixGid(unixGid)
			if err != nil {
				t.Errorf("Expected no error for %v, got %v", unixGid, err)
			}
		}
	})

	t.Run("InvalidUnixGid", func(t *testing.T) {
		invalidUnixGids := []int{-1, 1000000000000000000}
		for _, unixGid := range invalidUnixGids {
			_, err := NewUnixGid(unixGid)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", unixGid)
			}
		}
	})
}
