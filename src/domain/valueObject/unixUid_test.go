package valueObject

import "testing"

func TestUnixUid(t *testing.T) {
	t.Run("ValidUnixUid", func(t *testing.T) {
		validUnixUids := []int{0, 1000, 65365}
		for _, unixUid := range validUnixUids {
			_, err := NewUnixUid(unixUid)
			if err != nil {
				t.Errorf("Expected no error for %v, got %v", unixUid, err)
			}
		}
	})

	t.Run("InvalidUnixUid", func(t *testing.T) {
		invalidUnixUids := []int{-1, 1000000000000000000}
		for _, unixUid := range invalidUnixUids {
			_, err := NewUnixUid(unixUid)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", unixUid)
			}
		}
	})
}
