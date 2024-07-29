package valueObject

import (
	"testing"
)

func TestDatabasePrivilege(t *testing.T) {
	t.Run("ValidDatabasePrivilege", func(t *testing.T) {
		validDatabasePrivileges := []string{
			"ALL PRIVILEGES",
			"all",
			"ALTER ROUTINE",
			"alter system",
			"ALTER",
			"bypassrls",
		}

		for _, dbPrivilege := range validDatabasePrivileges {
			_, err := NewDatabasePrivilege(dbPrivilege)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", dbPrivilege, err)
			}
		}
	})

	t.Run("InvalidDatabasePrivilege", func(t *testing.T) {
		invalidDatabasePrivileges := []string{
			"-abc-123-xyz",
			"abc-123-",
			"ab",
			"a!b@c#123",
			"a-b-c-d-e-f-g-h-i-j-k-l-m-n-o-p-q",
		}

		for _, dbPrivilege := range invalidDatabasePrivileges {
			_, err := NewDatabasePrivilege(dbPrivilege)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dbPrivilege)
			}
		}
	})
}
