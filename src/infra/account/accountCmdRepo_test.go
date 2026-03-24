package accountInfra

import (
	"os"
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func addDummyUser() error {
	username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))
	password, _ := tkValueObject.NewPassword(os.Getenv("DUMMY_USER_PASS"))

	ipAddress := tkValueObject.IpAddressLocal
	operatorAccountId, _ := tkValueObject.NewAccountId(0)
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
	accountId, _ := tkValueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))
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
	accountId, _ := tkValueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))

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
		password, _ := tkValueObject.NewPassword("Invalid@123")
		ipAddress := tkValueObject.IpAddressLocal
		operatorAccountId, _ := tkValueObject.NewAccountId(0)
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

		newPassword, _ := tkValueObject.NewPassword("newPassword@1")
		updateDto := dto.NewUpdateAccount(
			&accountId, nil, &newPassword, nil, nil, accountId,
			tkValueObject.IpAddressLocal,
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
