package secureAccessKeyInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
)

func TestSecureAccessKeyCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	accountId, _ := valueObject.NewAccountId(2000)
	accountUsername, _ := valueObject.NewUsername("accountToTestSecureAccessKey")
	accountPassword, _ := valueObject.NewPassword("q1w2e3r4")
	ipAddress := valueObject.NewLocalhostIpAddress()

	createDto := dto.NewCreateAccount(
		accountUsername, accountPassword, accountId, ipAddress,
	)

	accountCmdRepo := accountInfra.NewAccountCmdRepo(testHelpers.GetPersistentDbSvc())
	_, err := accountCmdRepo.Create(createDto)
	if err != nil {
		t.Fatalf("FailToCreateTestAccount")
	}

	t.Skip("SkipSecureAccessKeysTests")

	var keyId valueObject.SecureAccessKeyId

	keyName, _ := valueObject.NewSecureAccessKeyName("testSecureAccessKey")
	secureAccessKeyCmdRepo := NewSecureAccessKeyCmdRepo(testHelpers.GetPersistentDbSvc())

	t.Run("CreateSecureAccessKey", func(t *testing.T) {
		keyContent, _ := valueObject.NewSecureAccessKeyContent(
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+GDqLA2sGauzU5hUxBbBmm6FfeZpUbiX6IlQO9KqeqAsum+Efhvj+qpatM5PzMMwtlcFwDS5Y4RcX9uxE8IGsYiALRfnLAX5p73zrcrXamMJSx25rXAu/VJdmekxHbDgsBPyk6/4dfu+3uW7ka7HHhPytPIqW2qBuPkalJinc7qKEuXdkCyX8+8a+0uN8XodLipLJwU8A1VPvI9thYxITyHWZnXRnin0r/unHgLrg9bBILXZf0JRslelYdCvuCGnRKZfokh153shMZ63S+iV/Tohg2bOVxyz3HIQ983ga24uTFQhLpITMe9JEfq3pp2wcCE5hNFlNKyeDG8kwB+8V",
		)
		createDto := dto.NewCreateSecureAccessKey(
			keyName, keyContent, accountId, accountId, ipAddress,
		)

		keyId, err = secureAccessKeyCmdRepo.Create(createDto)
		if err != nil {
			t.Fatalf(
				"Expected no error for %s, but got %s", keyName.String(), err.Error(),
			)
		}
	})

	t.Run("DeleteSecureAccessKey", func(t *testing.T) {
		deleteDto := dto.NewDeleteSecureAccessKey(keyId, accountId, accountId, ipAddress)

		err = secureAccessKeyCmdRepo.Delete(deleteDto)
		if err != nil {
			t.Fatalf(
				"Expected no error for %s, but got %s", keyName.String(), err.Error(),
			)
		}
	})

	err = accountCmdRepo.Delete(accountId)
	if err != nil {
		t.Fatalf("FailToDeleteTestAccount")
	}
}
