package valueObject

import (
	"testing"
)

func TestSecureAccessKeyFingerprint(t *testing.T) {
	t.Run("ValidSecureAccessKeyFingerprint", func(t *testing.T) {
		rawValidSecureAccessKeyFingerprint := []interface{}{
			"SHA256:+DZVNCZhuX6xKglL9R3mUkvRJpMeL8ptNi8kaxAShg4",
			"SHA256:4/A1a6zPZdue6c03mG9DBk7e0Mqt7167wK5ikSvxynw",
			"SHA256:fTmGqpEJy0oCGGobdzvH9KeBvPrQRFTxn1zr/ss4Wow",
		}

		for _, rawKeyFingerprint := range rawValidSecureAccessKeyFingerprint {
			_, err := NewSecureAccessKeyFingerprint(rawKeyFingerprint)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", rawKeyFingerprint, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSecureAccessKeyFingerprint", func(t *testing.T) {
		rawInvalidSecureAccessKeyFingerprint := []interface{}{
			"", "SHA256", ":+DZVNCZhuX6xKglL9R3mUkvRJpMeL8ptNi8kaxAShg4",
			"SHA256+DZVNCZhuX6xKglL9R3mUkvRJpMeL8ptNi8kaxAShg4",
			"+DZVNCZhuX6xKglL9R3mUkvRJpMeL8ptNi8kaxAShg4",
		}

		for _, rawKeyFingerprint := range rawInvalidSecureAccessKeyFingerprint {
			_, err := NewSecureAccessKeyFingerprint(rawKeyFingerprint)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", rawKeyFingerprint)
			}
		}
	})
}
