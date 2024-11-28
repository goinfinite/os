package accountInfra

import (
	"os"
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
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
}
