package valueObject

import "testing"

func TestAccessTokenType(t *testing.T) {
	t.Run("ValidAccessTokenType", func(t *testing.T) {
		validAccessTokenTypes := []interface{}{
			"sessionToken", "accountApiKey",
		}

		for _, accessTokenType := range validAccessTokenTypes {
			_, err := NewAccessTokenType(accessTokenType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", accessTokenType, err.Error())
			}
		}
	})

	t.Run("InvalidAccessTokenType", func(t *testing.T) {
		invalidAccessTokenTypes := []interface{}{
			"", "invalidAuthToken", "12345678",
		}

		for _, accessTokenType := range invalidAccessTokenTypes {
			_, err := NewAccessTokenType(accessTokenType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", accessTokenType)
			}
		}
	})
}
