package valueObject

import "testing"

func TestUnixCommand(t *testing.T) {
	t.Run("ValidUnixCommand", func(t *testing.T) {
		validUnixCommands := []interface{}{
			"curl https://google.com", "mv file1 file2", "os vhost get",
			"os services create-installable -n php",
		}
		for _, unixCommand := range validUnixCommands {
			_, err := NewUnixCommand(unixCommand)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", unixCommand, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidUnixCommand", func(t *testing.T) {
		invalidUnixCommands := []interface{}{
			"", "t",
		}
		for _, unixCommand := range invalidUnixCommands {
			_, err := NewUnixCommand(unixCommand)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", unixCommand)
			}
		}
	})
}
