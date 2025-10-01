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

	ipAddress := valueObject.IpAddressSystem
	operatorAccountId, _ := valueObject.NewAccountId(0)
	createDto := dto.NewCreateAccount(
		username, password, false, operatorAccountId, ipAddress,
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
		ipAddress := valueObject.IpAddressSystem
		operatorAccountId, _ := valueObject.NewAccountId(0)
		createDto := dto.NewCreateAccount(
			username, password, false, operatorAccountId, ipAddress,
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

	t.Run("UpdateValidAccount", func(t *testing.T) {
		resetDummyUser()

		newPassword, _ := valueObject.NewPassword("newPassword")
		updateDto := dto.NewUpdateAccount(
			&accountId, nil, &newPassword, nil, nil, accountId,
			valueObject.IpAddressSystem,
		)
		err := accountCmdRepo.Update(updateDto)
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
}
