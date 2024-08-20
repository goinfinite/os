package valueObject

import "testing"

func TestServiceStatus(t *testing.T) {
	t.Run("ValidServiceStatuses", func(t *testing.T) {
		validStatusAndAliases := []interface{}{
			"running", "run", "up", "true", "false", "off", "no", "stop", "stopped",
			"halt", "uninstall", "uninstalled", "remove", "purge",
		}

		for _, status := range validStatusAndAliases {
			_, err := NewServiceStatus(status)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", status, err.Error())
			}
		}
	})

	t.Run("InvalidServiceStatuses", func(t *testing.T) {
		invalidStatusAndAliases := []interface{}{
			"runningg", "runn", "upp", "truee", "falsee", "offf", "runn1ng",
			"un11install",
		}

		for _, status := range invalidStatusAndAliases {
			_, err := NewServiceStatus(status)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", status)
			}
		}
	})
}
