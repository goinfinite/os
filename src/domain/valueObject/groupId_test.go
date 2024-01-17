package valueObject

import "testing"

func TestGroupId(t *testing.T) {
	t.Run("ValidGroupId", func(t *testing.T) {
		validGroupIds := []interface{}{0, 1000, 65365, "12345"}
		for _, groupId := range validGroupIds {
			_, err := NewGroupId(groupId)
			if err != nil {
				t.Errorf("Expected no error for %v, got %v", groupId, err)
			}
		}
	})

	t.Run("InvalidGroupId", func(t *testing.T) {
		invalidGroupIds := []interface{}{-1, 1000000000000000000, "-455"}
		for _, groupId := range invalidGroupIds {
			_, err := NewGroupId(groupId)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", groupId)
			}
		}
	})
}
