package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func addDummyUser() error {
	username := valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME"))
	password := valueObject.NewPasswordPanic(os.Getenv("DUMMY_USER_PASS"))

	addUser := dto.AddAccount{
		Username: username,
		Password: password,
	}

	accCmdRepo := AccCmdRepo{}
	err := accCmdRepo.Add(addUser)
	if err != nil {
		return err
	}

	return nil
}

func deleteDummyUser() error {
	accountId := valueObject.NewAccountIdFromStringPanic(os.Getenv("DUMMY_USER_ID"))

	accCmdRepo := AccCmdRepo{}
	err := accCmdRepo.Delete(accountId)
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

func TestAccCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("AddValidAccount", func(t *testing.T) {
		err := addDummyUser()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("AddInvalidAccount", func(t *testing.T) {
		username := valueObject.NewUsernamePanic("root")
		password := valueObject.NewPasswordPanic("invalid")

		addUser := dto.AddAccount{
			Username: username,
			Password: password,
		}

		accCmdRepo := AccCmdRepo{}
		err := accCmdRepo.Add(addUser)
		if err == nil {
			t.Error("ExpectingError")
		}
	})

	t.Run("DeleteValidAccount", func(t *testing.T) {
		_ = addDummyUser()

		err := deleteDummyUser()
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		_ = addDummyUser()
	})

	t.Run("UpdatePasswordValidAccount", func(t *testing.T) {
		resetDummyUser()

		accountId := valueObject.NewAccountIdFromStringPanic(os.Getenv("DUMMY_USER_ID"))
		newPassword := valueObject.NewPasswordPanic("newPassword")

		accCmdRepo := AccCmdRepo{}
		err := accCmdRepo.UpdatePassword(accountId, newPassword)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		resetDummyUser()
	})

	t.Run("UpdateApiKeyValidAccount", func(t *testing.T) {
		resetDummyUser()

		accountId := valueObject.NewAccountIdFromStringPanic(os.Getenv("DUMMY_USER_ID"))

		accCmdRepo := AccCmdRepo{}
		_, err := accCmdRepo.UpdateApiKey(accountId)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		resetDummyUser()
	})
}
