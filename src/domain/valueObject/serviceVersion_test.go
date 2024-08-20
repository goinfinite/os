package valueObject

import "testing"

func TestServiceVersion(t *testing.T) {
	t.Run("ValidServiceVersions", func(t *testing.T) {
		validVersionsAndAliases := []interface{}{
			"1.0.0", "0.1.0", "latest", "lts", "alpha", "beta", "version1.0.0",
		}
		for _, name := range validVersionsAndAliases {
			_, err := NewServiceVersion(name)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", name, err.Error())
			}
		}
	})

	t.Run("InvalidServiceVersions", func(t *testing.T) {
		invalidVersionsAndAliases := []interface{}{
			"", "1.0<0",
		}
		for _, name := range invalidVersionsAndAliases {
			_, err := NewServiceVersion(name)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", name)
			}
		}
	})
}
