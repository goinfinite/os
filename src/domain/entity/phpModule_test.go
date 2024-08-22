package entity

import "testing"

func TestPhpModule(t *testing.T) {
	t.Run("ValidPhpModule", func(t *testing.T) {
		validPhpModules := []string{
			"curl:true", "opcache:TRUE", "mysqli:false", "apcu:FALSE",
		}

		for _, module := range validPhpModules {
			_, err := NewPhpModuleFromString(module)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", module, err.Error())
			}
		}
	})

	t.Run("InvalidPhpModule", func(t *testing.T) {
		invalidPhpModules := []string{
			"", "curl", "opcache:",
		}

		for _, module := range invalidPhpModules {
			_, err := NewPhpModuleFromString(module)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", module)
			}
		}
	})
}
