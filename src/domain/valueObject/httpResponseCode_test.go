package valueObject

import (
	"testing"
)

func TestNewHttpResponseCode(t *testing.T) {
	t.Run("ValidResponseCode", func(t *testing.T) {
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
				t.Errorf("ExpectingNoErrorButGot: %s", err.Error())
			}
		}
	})

	t.Run("InvalidResponseCode", func(t *testing.T) {
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
				t.Errorf("ExpectingErrorButDidNotGetFor: %v", responseCode)
			}
		}
	})
}
