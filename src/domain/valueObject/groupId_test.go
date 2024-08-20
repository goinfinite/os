package valueObject

import "testing"

func TestGroupId(t *testing.T) {
	t.Run("ValidGroupId", func(t *testing.T) {
		validGroupIds := []interface{}{
			0, 1, 10000000000000, "455", 40.5,
		}

		for _, groupId := range validGroupIds {
			_, err := NewGroupId(groupId)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", groupId, err.Error())
			}
		}
	})

	t.Run("InvalidGroupId", func(t *testing.T) {
		invalidGroupIds := []interface{}{
			-1, -10000000000000, "-455", -40.5,
		}

		for _, groupId := range invalidGroupIds {
			_, err := NewGroupId(groupId)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", groupId)
			}
		}
	})
}
