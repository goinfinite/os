package valueObject

import "testing"

func TestServiceName(t *testing.T) {
	t.Run("ValidServiceNames", func(t *testing.T) {
		validNamesAndAliases := []string{
			"openlitespeed",
			"litespeed",
			"nginx",
			"node",
			"nodejs",
			"redis-server",
		}
		for _, name := range validNamesAndAliases {
			_, err := NewServiceName(name)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", name, err)
			}
		}
	})

	t.Run("InvalidServiceNames", func(t *testing.T) {
		invalidNamesAndAliases := []string{
			"nginx@",
			"my<>sql",
			"php#fpm",
			"node(js)",
		}
		for _, name := range invalidNamesAndAliases {
			_, err := NewServiceName(name)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", name)
			}
		}
	})
}
