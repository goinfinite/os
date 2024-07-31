package valueObject

import "testing"

func TestDatabaseName(t *testing.T) {
	t.Run("ValidDatabaseName", func(t *testing.T) {
		validDatabaseNames := []string{
			"abc-123-xyz",
			"a1-b2-c3-d4-e5",
			"username-1234",
			"a-b-c-d-e-f-g-h-i-j-k-l",
			"valid-value-12345",
		}

		for _, dbName := range validDatabaseNames {
			_, err := NewDatabaseName(dbName)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", dbName, err)
			}
		}
	})

	t.Run("InvalidDatabaseName", func(t *testing.T) {
		invalidDatabaseNames := []string{
			"-abc-123-xyz",
			"abc-123-",
			"ab",
			"a!b@c#123",
			"a-b-c-d-e-f-g-h-i-j-k-l-m-n-o-p-q",
		}

		for _, dbName := range invalidDatabaseNames {
			_, err := NewDatabaseName(dbName)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dbName)
			}
		}
	})
}
