package valueObject

import (
	"testing"
)

func TestSecureAccessPublicKeyFingerprint(t *testing.T) {
	t.Run("ValidSecureAccessPublicKeyFingerprint", func(t *testing.T) {
		rawValidSecureAccessPublicKeyFingerprint := []interface{}{
			"SHA256:+DZVNCZhuX6xKglL9R3mUkvRJpMeL8ptNi8kaxAShg4",
			"SHA256:4/A1a6zPZdue6c03mG9DBk7e0Mqt7167wK5ikSvxynw",
			"SHA256:fTmGqpEJy0oCGGobdzvH9KeBvPrQRFTxn1zr/ss4Wow",
		}

		for _, rawKeyFingerprint := range rawValidSecureAccessPublicKeyFingerprint {
			_, err := NewSecureAccessPublicKeyFingerprint(rawKeyFingerprint)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", rawKeyFingerprint, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSecureAccessPublicKeyFingerprint", func(t *testing.T) {
		rawInvalidSecureAccessPublicKeyFingerprint := []interface{}{
			"", "SHA256", ":+DZVNCZhuX6xKglL9R3mUkvRJpMeL8ptNi8kaxAShg4",
			"SHA256+DZVNCZhuX6xKglL9R3mUkvRJpMeL8ptNi8kaxAShg4",
			"+DZVNCZhuX6xKglL9R3mUkvRJpMeL8ptNi8kaxAShg4",
		}

		for _, rawKeyFingerprint := range rawInvalidSecureAccessPublicKeyFingerprint {
			_, err := NewSecureAccessPublicKeyFingerprint(rawKeyFingerprint)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", rawKeyFingerprint)
			}
		}
	})
}
