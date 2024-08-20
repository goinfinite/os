package valueObject

import "testing"

func TestPhpModuleName(t *testing.T) {
	t.Run("ValidPhpModuleNames", func(t *testing.T) {
		validNames := []interface{}{
			"ioncube", "apcu", "imagick", "opcache", "mysqli",
		}

		for _, name := range validNames {
			_, err := NewPhpModuleName(name)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", name, err.Error())
			}
		}
	})

	t.Run("InvalidPhpModuleNames", func(t *testing.T) {
		invalidNames := []interface{}{
			"ioncube_loader.so", "<script>alert('xss')</script>", "@blabla@",
		}

		for _, name := range invalidNames {
			_, err := NewPhpModuleName(name)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", name)
			}
		}
	})
}
