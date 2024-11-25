package accountInfra

import (
	"os"
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestAccountQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	accountQueryRepo := NewAccountQueryRepo(persistentDbSvc)

	accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))

	t.Run("ReadValidAccounts", func(t *testing.T) {
		_, err := accountQueryRepo.Read()
		if err != nil {
			t.Errorf("Expecting no error, but got %s", err.Error())
		}
	})

	t.Run("ReadValidAccountByUsername", func(t *testing.T) {
		username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))

		_, err := accountQueryRepo.ReadByUsername(username)
		if err != nil {
			t.Errorf(
				"Expecting no error for %s, but got %s", username.String(), err.Error(),
			)
		}
	})

	t.Run("ReadValidAccountById", func(t *testing.T) {
		_, err := accountQueryRepo.ReadById(accountId)
		if err != nil {
			t.Errorf(
				"Expecting no error for %d, but got %s", accountId.Uint64(),
				err.Error(),
			)
		}
	})

	t.Run("ReadInvalidAccount", func(t *testing.T) {
		username, _ := valueObject.NewUsername("invalid")

		_, err := accountQueryRepo.ReadByUsername(username)
		if err == nil {
			t.Errorf("Expecting error for %s, but got nil", username.String())
		}
	})

	t.Skip("SkipSecureAccessKeysTests")

	accountCmdRepo := NewAccountCmdRepo(persistentDbSvc)

	keyName, _ := valueObject.NewSecureAccessKeyName("dummySecureAccessKey")
	keyContent, _ := valueObject.NewSecureAccessKeyContent(
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+GDqLA2sGauzU5hUxBbBmm6FfeZpUbiX6IlQO9KqeqAsum+Efhvj+qpatM5PzMMwtlcFwDS5Y4RcX9uxE8IGsYiALRfnLAX5p73zrcrXamMJSx25rXAu/VJdmekxHbDgsBPyk6/4dfu+3uW7ka7HHhPytPIqW2qBuPkalJinc7qKEuXdkCyX8+8a+0uN8XodLipLJwU8A1VPvI9thYxITyHWZnXRnin0r/unHgLrg9bBILXZf0JRslelYdCvuCGnRKZfokh153shMZ63S+iV/Tohg2bOVxyz3HIQ983ga24uTFQhLpITMe9JEfq3pp2wcCE5hNFlNKyeDG8kwB+8V",
	)
	createDto := dto.NewCreateSecureAccessKey(
		keyName, keyContent, accountId, accountId, valueObject.NewLocalhostIpAddress(),
	)

	_, err := accountCmdRepo.CreateSecureAccessKey(createDto)
	if err != nil {
		t.Fatalf("Fail to create dummy SecureAccessKey to test")
	}

	t.Run("ReadSecureAccessKeys", func(t *testing.T) {
		keys, err := accountQueryRepo.ReadSecureAccessKeys(accountId)
		if err != nil {
			t.Fatalf(
				"Expecting no error for %d, but got %s", accountId.Uint64(), err.Error(),
			)
		}

		if len(keys) == 0 {
			t.Error("Expecting a keys list, but got an empty one")
		}
	})

	t.Run("ReadSecureAccessKeyByName", func(t *testing.T) {
		_, err := accountQueryRepo.ReadSecureAccessKeyByName(accountId, keyName)
		if err != nil {
			t.Fatalf(
				"Expecting no error for %s (%d), but got %s", keyName.String(),
				accountId.Uint64(), err.Error(),
			)
		}
	})
}
