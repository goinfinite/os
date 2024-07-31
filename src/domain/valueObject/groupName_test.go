package valueObject

import "testing"

func TestGroupName(t *testing.T) {
	t.Run("ValidGroupName", func(t *testing.T) {
		validGroupNames := []interface{}{
			"ssl-cert",
			"damn-man--",
			"root",
			"mysql",
		}

		for _, groupName := range validGroupNames {
			_, err := NewGroupName(groupName)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", groupName, err)
			}
		}
	})

	t.Run("InvalidGroupName", func(t *testing.T) {
		invalidGroupNames := []interface{}{
			"",
			".",
			"..",
			"/",
			"root:root",
			"roooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooot",
			"ççççççç",
			"<root>",
			"not a valid user",
		}

		for _, groupName := range invalidGroupNames {
			_, err := NewGroupName(groupName)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", groupName)
			}
		}
	})
}
