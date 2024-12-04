package valueObject

import "testing"

func TestSecureAccessPublicKeyContent(t *testing.T) {
	t.Run("ValidSecureAccessPublicKeyContent", func(t *testing.T) {
		rawValidSecureAccessPublicKeyContent := []interface{}{
			"ssh-rsa c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U= myMachine@pop-os",
			"ssh-rsa c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U=",
			"ssh-ed25519 c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U= myMachine@pop-os",
			"ssh-ed25519 c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U=",
		}

		for _, rawKeyContent := range rawValidSecureAccessPublicKeyContent {
			_, err := NewSecureAccessPublicKeyContent(rawKeyContent)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", rawKeyContent, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSecureAccessPublicKeyContent", func(t *testing.T) {
		rawInvalidSecureAccessPublicKeyContent := []interface{}{
			12345, 1.25, true, "", "ssh-rsa", "ssh-rsa myMachine@pop-os",
			"c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U= myMachine@pop-os",
		}

		for _, rawKeyContent := range rawInvalidSecureAccessPublicKeyContent {
			_, err := NewSecureAccessPublicKeyContent(rawKeyContent)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", rawKeyContent)
			}
		}
	})

	t.Run("ReadOnlyKeyName", func(t *testing.T) {
		rawValidSecureAccessPublicKeyContentWithName := "ssh-rsa c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U= myMachine@pop-os"
		keyContentWithName, err := NewSecureAccessPublicKeyContent(rawValidSecureAccessPublicKeyContentWithName)
		if err != nil {
			t.Fatalf(
				"Expected no error for '%v', got '%s'",
				rawValidSecureAccessPublicKeyContentWithName, err.Error(),
			)
		}

		_, err = keyContentWithName.ReadOnlyKeyName()
		if err != nil {
			t.Errorf(
				"Expected no error for '%v', got '%s'",
				rawValidSecureAccessPublicKeyContentWithName, err.Error(),
			)
		}

		rawValidSecureAccessPublicKeyContentWithoutName := "ssh-ed25519 c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U="
		keyContentWithoutName, err := NewSecureAccessPublicKeyContent(rawValidSecureAccessPublicKeyContentWithoutName)
		if err != nil {
			t.Fatalf(
				"Expected no error for '%v', got '%s'",
				rawValidSecureAccessPublicKeyContentWithName, err.Error(),
			)
		}

		_, err = keyContentWithoutName.ReadOnlyKeyName()
		if err == nil {
			t.Errorf(
				"Expected error for '%v', got nil",
				rawValidSecureAccessPublicKeyContentWithName,
			)
		}
	})
}
