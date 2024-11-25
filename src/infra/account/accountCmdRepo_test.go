package accountInfra

import (
	"os"
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func addDummyUser() error {
	username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))
	password, _ := valueObject.NewPassword(os.Getenv("DUMMY_USER_PASS"))

	ipAddress := valueObject.NewLocalhostIpAddress()
	operatorAccountId, _ := valueObject.NewAccountId(0)
	createDto := dto.NewCreateAccount(
		username, password, operatorAccountId, ipAddress,
	)

	accountCmdRepo := NewAccountCmdRepo(testHelpers.GetPersistentDbSvc())
	_, err := accountCmdRepo.Create(createDto)
	if err != nil {
		return err
	}

	return nil
}

func deleteDummyUser() error {
	accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))
	accountCmdRepo := NewAccountCmdRepo(testHelpers.GetPersistentDbSvc())
	err := accountCmdRepo.Delete(accountId)
	if err != nil {
		return err
	}

	return nil
}

func resetDummyUser() {
	_ = addDummyUser()
	_ = deleteDummyUser()
	_ = addDummyUser()
}

func TestAccountCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	accountCmdRepo := NewAccountCmdRepo(testHelpers.GetPersistentDbSvc())
	accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))

	t.Run("AddValidAccount", func(t *testing.T) {
		err := addDummyUser()
		if err != nil {
			t.Errorf(
				"Expected no error for %d, but got %s", accountId.Uint64(), err.Error(),
			)
		}
	})

	t.Run("AddInvalidAccount", func(t *testing.T) {
		username, _ := valueObject.NewUsername("root")
		password, _ := valueObject.NewPassword("invalid")
		ipAddress := valueObject.NewLocalhostIpAddress()
		operatorAccountId, _ := valueObject.NewAccountId(0)
		createDto := dto.NewCreateAccount(
			username, password, operatorAccountId, ipAddress,
		)

		_, err := accountCmdRepo.Create(createDto)
		if err == nil {
			t.Error("AccountShouldNotBeAdded")
		}
	})

	t.Run("DeleteValidAccount", func(t *testing.T) {
		err := deleteDummyUser()
		if err != nil {
			t.Errorf(
				"Expected no error for %d, but got %s", accountId.Uint64(), err.Error(),
			)
		}
	})

	t.Run("UpdatePasswordValidAccount", func(t *testing.T) {
		resetDummyUser()

		newPassword, _ := valueObject.NewPassword("newPassword")

		err := accountCmdRepo.UpdatePassword(accountId, newPassword)
		if err != nil {
			t.Errorf(
				"Expected no error for %s, but got %s", newPassword.String(),
				err.Error(),
			)
		}
	})

	t.Run("UpdateApiKeyValidAccount", func(t *testing.T) {
		resetDummyUser()

		_, err := accountCmdRepo.UpdateApiKey(accountId)
		if err != nil {
			t.Errorf(
				"Expected no error for %d, but got %s", accountId.Uint64(), err.Error(),
			)
		}
	})

	t.Skip("SkipSecureAccessKeysTests")

	t.Run("CreateSecureAccessKey", func(t *testing.T) {
		keyName, _ := valueObject.NewSecureAccessKeyName("dummySecureAccessKey")
		keyContent, _ := valueObject.NewSecureAccessKeyContent(
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+GDqLA2sGauzU5hUxBbBmm6FfeZpUbiX6IlQO9KqeqAsum+Efhvj+qpatM5PzMMwtlcFwDS5Y4RcX9uxE8IGsYiALRfnLAX5p73zrcrXamMJSx25rXAu/VJdmekxHbDgsBPyk6/4dfu+3uW7ka7HHhPytPIqW2qBuPkalJinc7qKEuXdkCyX8+8a+0uN8XodLipLJwU8A1VPvI9thYxITyHWZnXRnin0r/unHgLrg9bBILXZf0JRslelYdCvuCGnRKZfokh153shMZ63S+iV/Tohg2bOVxyz3HIQ983ga24uTFQhLpITMe9JEfq3pp2wcCE5hNFlNKyeDG8kwB+8V",
		)
		createDto := dto.NewCreateSecureAccessKey(
			keyName, keyContent, accountId, accountId,
			valueObject.NewLocalhostIpAddress(),
		)

		_, err := accountCmdRepo.CreateSecureAccessKey(createDto)
		if err != nil {
			t.Fatalf(
				"Expected no error for %s, but got %s", keyName.String(), err.Error(),
			)
		}
	})
}
