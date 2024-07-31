package valueObject

import "testing"

func TestAccessTokenStr(t *testing.T) {
	t.Run("ValidAccessTokenStr", func(t *testing.T) {
		validAccessTokenStrs := []interface{}{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODEyMzkwIiwibmFtZSI6IlRoYXQgbWUiLCJpYXQiOjE1MTYyMzM0OTAyMn0.GopUUcCOZIrmTipkcmDUsOrD0m5ymDfK-p3wHWNYwQk",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3MTIzMTIzODEyMzkwIiwibmFtZSI6IlRoYXQgbWFhc2RhZHNhZGUiLCJpYXQiOjE1MTYyMzM0OTAyMn0.4bT9W57v1TkHi7Ern_2iRgmbHW4jBr4IPuvs9_dQNww",
		}

		for _, accessTokenStr := range validAccessTokenStrs {
			_, err := NewAccessTokenStr(accessTokenStr)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", accessTokenStr, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidAccessTokenStr", func(t *testing.T) {
		invalidAccessTokenStrs := []interface{}{
			"", "invalidAuthToken", "12345678",
		}

		for _, accessTokenStr := range invalidAccessTokenStrs {
			_, err := NewAccessTokenStr(accessTokenStr)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", accessTokenStr)
			}
		}
	})
}
