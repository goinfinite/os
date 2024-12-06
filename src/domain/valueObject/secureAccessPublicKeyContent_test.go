package valueObject

import "testing"

func TestSecureAccessPublicKeyContent(t *testing.T) {
	t.Run("ValidSecureAccessPublicKeyContent", func(t *testing.T) {
		rawValidSecureAccessPublicKeyContent := []interface{}{
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCvDkVs/zS9pDcKY+0EC6koQDcb4yINUohwYr97JXEvuaFXVf8lUWnPn9nK9R9pH3srbuFUrFkJam/DvGgOlJIcH0uuWlj/ffQOs1Ts3rV86MD29XV2/bA9gLJD6rLCR1WJIKmnjozFcgaB3AwOH7/YkENKXJcBfO4WRHMcZzzrjUsVTsBSO3+EDxBfPmpjXKHyTkdCQ3EC48tV01lyXe7IWLUKCc3nN5Hv14/fA+lvtiTvR4WpXXiHrXFxh9xy381FdVZxQ6xYfjE+SbI1h7XHvaDQo6lglZFuqFftQtuo/QmNz3OLCc/oGNw202igxx8Iv/NBJLEr+6DRDwhDzO39RUQ7mRqr5coIcnf1uYZgCLUnq6md9sEll6OpsCSHDnCgi1LLrOa4ZnC/JGCfHO4yAbZxw7Yc3u9jP29d9zlGoTBx+G60JBIeGKGKdMYOAfQGDZp1uwiwdIS0aM15ph6c0/6mdrQw8ynSVqF5o+uh8FHYXC4DgIGAmtZR7Nna4+U= myMachine@pop-os",
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCvDkVs/zS9pDcKY+0EC6koQDcb4yINUohwYr97JXEvuaFXVf8lUWnPn9nK9R9pH3srbuFUrFkJam/DvGgOlJIcH0uuWlj/ffQOs1Ts3rV86MD29XV2/bA9gLJD6rLCR1WJIKmnjozFcgaB3AwOH7/YkENKXJcBfO4WRHMcZzzrjUsVTsBSO3+EDxBfPmpjXKHyTkdCQ3EC48tV01lyXe7IWLUKCc3nN5Hv14/fA+lvtiTvR4WpXXiHrXFxh9xy381FdVZxQ6xYfjE+SbI1h7XHvaDQo6lglZFuqFftQtuo/QmNz3OLCc/oGNw202igxx8Iv/NBJLEr+6DRDwhDzO39RUQ7mRqr5coIcnf1uYZgCLUnq6md9sEll6OpsCSHDnCgi1LLrOa4ZnC/JGCfHO4yAbZxw7Yc3u9jP29d9zlGoTBx+G60JBIeGKGKdMYOAfQGDZp1uwiwdIS0aM15ph6c0/6mdrQw8ynSVqF5o+uh8FHYXC4DgIGAmtZR7Nna4+U=",
			"ssh-ed25519 AAAAB3NzaC1yc2EAAAADAQABAAABgQCvDkVs/zS9pDcKY+0EC6koQDcb4yINUohwYr97JXEvuaFXVf8lUWnPn9nK9R9pH3srbuFUrFkJam/DvGgOlJIcH0uuWlj/ffQOs1Ts3rV86MD29XV2/bA9gLJD6rLCR1WJIKmnjozFcgaB3AwOH7/YkENKXJcBfO4WRHMcZzzrjUsVTsBSO3+EDxBfPmpjXKHyTkdCQ3EC48tV01lyXe7IWLUKCc3nN5Hv14/fA+lvtiTvR4WpXXiHrXFxh9xy381FdVZxQ6xYfjE+SbI1h7XHvaDQo6lglZFuqFftQtuo/QmNz3OLCc/oGNw202igxx8Iv/NBJLEr+6DRDwhDzO39RUQ7mRqr5coIcnf1uYZgCLUnq6md9sEll6OpsCSHDnCgi1LLrOa4ZnC/JGCfHO4yAbZxw7Yc3u9jP29d9zlGoTBx+G60JBIeGKGKdMYOAfQGDZp1uwiwdIS0aM15ph6c0/6mdrQw8ynSVqF5o+uh8FHYXC4DgIGAmtZR7Nna4+U= myMachine@pop-os",
			"ssh-ed25519 AAAAB3NzaC1yc2EAAAADAQABAAABgQCvDkVs/zS9pDcKY+0EC6koQDcb4yINUohwYr97JXEvuaFXVf8lUWnPn9nK9R9pH3srbuFUrFkJam/DvGgOlJIcH0uuWlj/ffQOs1Ts3rV86MD29XV2/bA9gLJD6rLCR1WJIKmnjozFcgaB3AwOH7/YkENKXJcBfO4WRHMcZzzrjUsVTsBSO3+EDxBfPmpjXKHyTkdCQ3EC48tV01lyXe7IWLUKCc3nN5Hv14/fA+lvtiTvR4WpXXiHrXFxh9xy381FdVZxQ6xYfjE+SbI1h7XHvaDQo6lglZFuqFftQtuo/QmNz3OLCc/oGNw202igxx8Iv/NBJLEr+6DRDwhDzO39RUQ7mRqr5coIcnf1uYZgCLUnq6md9sEll6OpsCSHDnCgi1LLrOa4ZnC/JGCfHO4yAbZxw7Yc3u9jP29d9zlGoTBx+G60JBIeGKGKdMYOAfQGDZp1uwiwdIS0aM15ph6c0/6mdrQw8ynSVqF5o+uh8FHYXC4DgIGAmtZR7Nna4+U=",
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
		rawValidSecureAccessPublicKeyContentWithName := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCvDkVs/zS9pDcKY+0EC6koQDcb4yINUohwYr97JXEvuaFXVf8lUWnPn9nK9R9pH3srbuFUrFkJam/DvGgOlJIcH0uuWlj/ffQOs1Ts3rV86MD29XV2/bA9gLJD6rLCR1WJIKmnjozFcgaB3AwOH7/YkENKXJcBfO4WRHMcZzzrjUsVTsBSO3+EDxBfPmpjXKHyTkdCQ3EC48tV01lyXe7IWLUKCc3nN5Hv14/fA+lvtiTvR4WpXXiHrXFxh9xy381FdVZxQ6xYfjE+SbI1h7XHvaDQo6lglZFuqFftQtuo/QmNz3OLCc/oGNw202igxx8Iv/NBJLEr+6DRDwhDzO39RUQ7mRqr5coIcnf1uYZgCLUnq6md9sEll6OpsCSHDnCgi1LLrOa4ZnC/JGCfHO4yAbZxw7Yc3u9jP29d9zlGoTBx+G60JBIeGKGKdMYOAfQGDZp1uwiwdIS0aM15ph6c0/6mdrQw8ynSVqF5o+uh8FHYXC4DgIGAmtZR7Nna4+U= myMachine@pop-os"
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

		rawValidSecureAccessPublicKeyContentWithoutName := "ssh-ed25519 AAAAB3NzaC1yc2EAAAADAQABAAABgQCvDkVs/zS9pDcKY+0EC6koQDcb4yINUohwYr97JXEvuaFXVf8lUWnPn9nK9R9pH3srbuFUrFkJam/DvGgOlJIcH0uuWlj/ffQOs1Ts3rV86MD29XV2/bA9gLJD6rLCR1WJIKmnjozFcgaB3AwOH7/YkENKXJcBfO4WRHMcZzzrjUsVTsBSO3+EDxBfPmpjXKHyTkdCQ3EC48tV01lyXe7IWLUKCc3nN5Hv14/fA+lvtiTvR4WpXXiHrXFxh9xy381FdVZxQ6xYfjE+SbI1h7XHvaDQo6lglZFuqFftQtuo/QmNz3OLCc/oGNw202igxx8Iv/NBJLEr+6DRDwhDzO39RUQ7mRqr5coIcnf1uYZgCLUnq6md9sEll6OpsCSHDnCgi1LLrOa4ZnC/JGCfHO4yAbZxw7Yc3u9jP29d9zlGoTBx+G60JBIeGKGKdMYOAfQGDZp1uwiwdIS0aM15ph6c0/6mdrQw8ynSVqF5o+uh8FHYXC4DgIGAmtZR7Nna4+U="
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
