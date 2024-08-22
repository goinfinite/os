package valueObject

import "testing"

func TestUnixUid(t *testing.T) {
	t.Run("ValidUnixUid", func(t *testing.T) {
		validUnixUids := []interface{}{
			-10000, 0, 1, "455", 65365,
		}

		for _, unixUid := range validUnixUids {
			_, err := NewUnixUid(unixUid)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", unixUid, err.Error())
			}
		}
	})
}
