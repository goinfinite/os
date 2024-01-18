package valueObject

import "testing"

func TestUnixUid(t *testing.T) {
	t.Run("ValidUnixUid", func(t *testing.T) {
		validUnixUids := []interface{}{0, 1000, 65365, "12345"}
		for _, groupId := range validUnixUids {
			_, err := NewGroupId(groupId)
			if err != nil {
				t.Errorf("Expected no error for %v, got %v", groupId, err)
			}
		}
	})

	t.Run("InvalidUnixUid", func(t *testing.T) {
		invalidUnixUids := []interface{}{-1, 1000000000000000000, "-455"}
		for _, groupId := range invalidUnixUids {
			_, err := NewGroupId(groupId)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", groupId)
			}
		}
	})
}
