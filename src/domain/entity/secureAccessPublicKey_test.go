package entity

import (
	"os"
	"testing"

	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestSecureAccessPublicKey(t *testing.T) {
	publicKeyId, _ := valueObject.NewSecureAccessPublicKeyId(1)
	accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))
	publicKeyName, _ := valueObject.NewSecureAccessPublicKeyName("myMachine@pop-os")

	t.Run("ValidSecureAccessPublicKey (RSA)", func(t *testing.T) {
		publicKeyContent, _ := valueObject.NewSecureAccessPublicKeyContent(
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCvDkVs/zS9pDcKY+0EC6koQDcb4yINUohwYr97JXEvuaFXVf8lUWnPn9nK9R9pH3srbuFUrFkJam/DvGgOlJIcH0uuWlj/ffQOs1Ts3rV86MD29XV2/bA9gLJD6rLCR1WJIKmnjozFcgaB3AwOH7/YkENKXJcBfO4WRHMcZzzrjUsVTsBSO3+EDxBfPmpjXKHyTkdCQ3EC48tV01lyXe7IWLUKCc3nN5Hv14/fA+lvtiTvR4WpXXiHrXFxh9xy381FdVZxQ6xYfjE+SbI1h7XHvaDQo6lglZFuqFftQtuo/QmNz3OLCc/oGNw202igxx8Iv/NBJLEr+6DRDwhDzO39RUQ7mRqr5coIcnf1uYZgCLUnq6md9sEll6OpsCSHDnCgi1LLrOa4ZnC/JGCfHO4yAbZxw7Yc3u9jP29d9zlGoTBx+G60JBIeGKGKdMYOAfQGDZp1uwiwdIS0aM15ph6c0/6mdrQw8ynSVqF5o+uh8FHYXC4DgIGAmtZR7Nna4+U=",
		)

		_, err := NewSecureAccessPublicKey(
			publicKeyId, accountId, publicKeyContent, publicKeyName,
			valueObject.NewUnixTimeNow(), valueObject.NewUnixTimeNow(),
		)
		if err != nil {
			t.Errorf(
				"Expected no error for %s, but got %s", publicKeyContent.String(),
				err.Error(),
			)
		}
	})

	t.Run("InvalidecureAccessPublicKey (RSA)", func(t *testing.T) {
		publicKeyContent, _ := valueObject.NewSecureAccessPublicKeyContent(
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCvDkVs/zS9pDcKY+0EC6koQD+U=",
		)

		_, err := NewSecureAccessPublicKey(
			publicKeyId, accountId, publicKeyContent, publicKeyName,
			valueObject.NewUnixTimeNow(), valueObject.NewUnixTimeNow(),
		)
		if err == nil {
			t.Errorf("Expected error for %s, but got nil", publicKeyContent.String())
		}
	})

	t.Run("ValidSecureAccessPublicKey (Ed25519)", func(t *testing.T) {
		publicKeyContent, _ := valueObject.NewSecureAccessPublicKeyContent(
			"ssh-ed25519 AAAAB3NzaC1yc2EAAAADAQABAAABgQCvDkVs/zS9pDcKY+0EC6koQDcb4yINUohwYr97JXEvuaFXVf8lUWnPn9nK9R9pH3srbuFUrFkJam/DvGgOlJIcH0uuWlj/ffQOs1Ts3rV86MD29XV2/bA9gLJD6rLCR1WJIKmnjozFcgaB3AwOH7/YkENKXJcBfO4WRHMcZzzrjUsVTsBSO3+EDxBfPmpjXKHyTkdCQ3EC48tV01lyXe7IWLUKCc3nN5Hv14/fA+lvtiTvR4WpXXiHrXFxh9xy381FdVZxQ6xYfjE+SbI1h7XHvaDQo6lglZFuqFftQtuo/QmNz3OLCc/oGNw202igxx8Iv/NBJLEr+6DRDwhDzO39RUQ7mRqr5coIcnf1uYZgCLUnq6md9sEll6OpsCSHDnCgi1LLrOa4ZnC/JGCfHO4yAbZxw7Yc3u9jP29d9zlGoTBx+G60JBIeGKGKdMYOAfQGDZp1uwiwdIS0aM15ph6c0/6mdrQw8ynSVqF5o+uh8FHYXC4DgIGAmtZR7Nna4+U=",
		)

		_, err := NewSecureAccessPublicKey(
			publicKeyId, accountId, publicKeyContent, publicKeyName,
			valueObject.NewUnixTimeNow(), valueObject.NewUnixTimeNow(),
		)
		if err != nil {
			t.Errorf(
				"Expected no error for %s, but got %s", publicKeyContent.String(),
				err.Error(),
			)
		}
	})

	t.Run("InvalidecureAccessPublicKey (Ed25519)", func(t *testing.T) {
		publicKeyContent, _ := valueObject.NewSecureAccessPublicKeyContent(
			"ssh-ed25519 AAAAB3NzaC1yc2EAAAADAQABAAABgQCvDkVs/zS9pDcKY+0EC6koQD+U=",
		)

		_, err := NewSecureAccessPublicKey(
			publicKeyId, accountId, publicKeyContent, publicKeyName,
			valueObject.NewUnixTimeNow(), valueObject.NewUnixTimeNow(),
		)
		if err == nil {
			t.Errorf("Expected error for %s, but got nil", publicKeyContent.String())
		}
	})
}
