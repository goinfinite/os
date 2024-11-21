package valueObject

import "testing"

func TestSecureAccessKeyContent(t *testing.T) {
	t.Run("ValidSecureAccessKeyContent", func(t *testing.T) {
		rawValidSecureAccessKeyContent := []interface{}{
			"ssh-rsa c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U= myMachine@pop-os",
			"ssh-rsa c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U=",
		}

		for _, rawKeyContent := range rawValidSecureAccessKeyContent {
			_, err := NewSecureAccessKeyContent(rawKeyContent)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", rawKeyContent, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSecureAccessKeyContent", func(t *testing.T) {
		rawInvalidSecureAccessKeyContent := []interface{}{
			12345, 1.25, true, "", "ssh-rsa", "ssh-rsa myMachine@pop-os",
			"c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U= myMachine@pop-os",
		}

		for _, rawKeyContent := range rawInvalidSecureAccessKeyContent {
			_, err := NewSecureAccessKeyContent(rawKeyContent)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", rawKeyContent)
			}
		}
	})

	t.Run("ReadOnlyKeyName", func(t *testing.T) {
		rawValidSecureAccessKeyContentWithName := "ssh-rsa c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U= myMachine@pop-os"
		keyContentWithName, err := NewSecureAccessKeyContent(rawValidSecureAccessKeyContentWithName)
		if err != nil {
			t.Fatalf(
				"Expected no error for '%v', got '%s'",
				rawValidSecureAccessKeyContentWithName, err.Error(),
			)
		}

		_, err = keyContentWithName.ReadOnlyKeyName()
		if err != nil {
			t.Errorf(
				"Expected no error for '%v', got '%s'",
				rawValidSecureAccessKeyContentWithName, err.Error(),
			)
		}

		rawValidSecureAccessKeyContentWithoutName := "ssh-rsa c2VjdXJlIGFjY2Vzcy/BrZXkgY29udGV/udCB0ZXN0+U="
		keyContentWithoutName, err := NewSecureAccessKeyContent(rawValidSecureAccessKeyContentWithoutName)
		if err != nil {
			t.Fatalf(
				"Expected no error for '%v', got '%s'",
				rawValidSecureAccessKeyContentWithName, err.Error(),
			)
		}

		_, err = keyContentWithoutName.ReadOnlyKeyName()
		if err == nil {
			t.Errorf(
				"Expected error for '%v', got nil",
				rawValidSecureAccessKeyContentWithName,
			)
		}
	})
}
