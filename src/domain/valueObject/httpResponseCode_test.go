package valueObject

import "testing"

func TestNewHttpResponseCode(t *testing.T) {
	t.Run("ValidHttpResponseCode", func(t *testing.T) {
		validResponseCodes := []interface{}{
			"100",
			"200",
			"300",
			"400",
			"500",
			100,
			200,
			300,
			400,
			500,
		}

		for _, responseCode := range validResponseCodes {
			_, err := NewHttpResponseCode(responseCode)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v,' got '%s'",
					responseCode,
					err.Error(),
				)
			}
		}
	})

	t.Run("InvalidHttpResponseCode", func(t *testing.T) {
		invalidResponseCodes := []interface{}{
			"@blabla",
			"<script>alert('xss')</script>",
			"1000",
			"0",
			"-1",
			"UNION SELECT * FROM USERS",
		}

		for _, responseCode := range invalidResponseCodes {
			_, err := NewHttpResponseCode(responseCode)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", responseCode)
			}
		}
	})
}
