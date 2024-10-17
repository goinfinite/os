package valueObject

import "testing"

func TestActivityRecordCode(t *testing.T) {
	t.Run("ValidActivityRecordCode", func(t *testing.T) {
		validActivityRecordCodes := []interface{}{
			"LoginFailed", "LoginSuccessful", "AccountCreated", "AccountDeleted",
			"AccountPasswordUpdated", "AccountApiKeyUpdated", "AccountQuotaUpdated",
			"UnauthorizedAccess",
		}

		for _, code := range validActivityRecordCodes {
			_, err := NewActivityRecordCode(code)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", code, err.Error())
			}
		}
	})

	t.Run("InvalidActivityRecordCode", func(t *testing.T) {
		invalidActivityRecordCodes := []interface{}{
			"", "a", 1000,
		}

		for _, code := range invalidActivityRecordCodes {
			_, err := NewActivityRecordCode(code)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", code)
			}
		}
	})
}
