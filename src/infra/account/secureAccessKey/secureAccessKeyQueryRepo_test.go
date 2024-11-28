package secureAccessKeyInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
)

func TestSecureAccessKeyQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	accountId, _ := valueObject.NewAccountId(2000)
	accountUsername, _ := valueObject.NewUsername("accountToTestSecureAccessKey")
	accountPassword, _ := valueObject.NewPassword("q1w2e3r4")
	ipAddress := valueObject.NewLocalhostIpAddress()

	createAccountDto := dto.NewCreateAccount(
		accountUsername, accountPassword, accountId, ipAddress,
	)

	accountCmdRepo := accountInfra.NewAccountCmdRepo(testHelpers.GetPersistentDbSvc())
	_, err := accountCmdRepo.Create(createAccountDto)
	if err != nil {
		t.Fatalf("FailToCreateTestAccount")
	}

	t.Skip("SkipSecureAccessKeysTests")

	keyName, _ := valueObject.NewSecureAccessKeyName("testSecureAccessKey")
	keyContent, _ := valueObject.NewSecureAccessKeyContent(
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+GDqLA2sGauzU5hUxBbBmm6FfeZpUbiX6IlQO9KqeqAsum+Efhvj+qpatM5PzMMwtlcFwDS5Y4RcX9uxE8IGsYiALRfnLAX5p73zrcrXamMJSx25rXAu/VJdmekxHbDgsBPyk6/4dfu+3uW7ka7HHhPytPIqW2qBuPkalJinc7qKEuXdkCyX8+8a+0uN8XodLipLJwU8A1VPvI9thYxITyHWZnXRnin0r/unHgLrg9bBILXZf0JRslelYdCvuCGnRKZfokh153shMZ63S+iV/Tohg2bOVxyz3HIQ983ga24uTFQhLpITMe9JEfq3pp2wcCE5hNFlNKyeDG8kwB+8V",
	)
	createDto := dto.NewCreateSecureAccessKey(
		keyName, keyContent, accountId, accountId, ipAddress,
	)

	secureAccessKeyCmdRepo := NewSecureAccessKeyCmdRepo(testHelpers.GetPersistentDbSvc())
	keyId, err := secureAccessKeyCmdRepo.Create(createDto)
	if err != nil {
		t.Fatalf("Fail to create dummy SecureAccessKey to test")
	}

	secureAccessKeyQueryRepo := NewSecureAccessKeyQueryRepo(testHelpers.GetPersistentDbSvc())

	requestDto := dto.ReadSecureAccessKeysRequest{
		AccountId:         accountId,
		SecureAccessKeyId: &keyId,
	}
	t.Run("ReadSecureAccessKeys", func(t *testing.T) {
		responseDto, err := secureAccessKeyQueryRepo.Read(requestDto)
		if err != nil {
			t.Fatalf(
				"Expecting no error for %d, but got %s", accountId.Uint64(),
				err.Error(),
			)
		}

		if len(responseDto.SecureAccessKeys) == 0 {
			t.Error("Expecting a keys list, but got an empty one")
		}
	})

	t.Run("ReadSecureAccessKeyByName", func(t *testing.T) {
		_, err := secureAccessKeyQueryRepo.ReadFirst(requestDto)
		if err != nil {
			t.Fatalf(
				"Expecting no error for %s (%d), but got %s", keyName.String(),
				accountId.Uint64(), err.Error(),
			)
		}
	})

	err = accountCmdRepo.Delete(accountId)
	if err != nil {
		t.Fatalf("FailToDeleteTestAccount")
	}
}
