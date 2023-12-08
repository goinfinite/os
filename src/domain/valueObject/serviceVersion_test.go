package valueObject

import "testing"

func TestServiceVersion(t *testing.T) {
	t.Run("ValidServiceVersions", func(t *testing.T) {
		validVersionsAndAliases := []string{
			"1.0.0",
			"0.1.0",
			"latest",
			"lts",
			"alpha",
			"beta",
		}
		for _, name := range validVersionsAndAliases {
			_, err := NewServiceVersion(name)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", name, err)
			}
		}
	})

	t.Run("InvalidServiceVersions", func(t *testing.T) {
		invalidVersionsAndAliases := []string{
			"",
			"1.0<0",
			"version1.0.0",
		}
		for _, name := range invalidVersionsAndAliases {
			_, err := NewServiceVersion(name)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", name)
			}
		}
	})
}
