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

	id, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))

	addDummyUser()

	t.Run("ReadValid", func(t *testing.T) {
		requestDto := dto.ReadAccountsRequest{
			AccountId: &id,
		}
		_, err := accountQueryRepo.Read(requestDto)
		if err != nil {
			t.Errorf("Expecting no error, but got %s", err.Error())
		}
	})

	t.Run("ReadInvalid", func(t *testing.T) {
		username, _ := valueObject.NewUsername("invalid")
		requestDto := dto.ReadAccountsRequest{
			AccountId: &id,
		}

		_, err := accountQueryRepo.ReadFirst(requestDto)
		if err == nil {
			t.Errorf("Expecting error for %s, but got nil", username.String())
		}
	})

	t.Run("ReadFirstValid", func(t *testing.T) {
		username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))
		requestDto := dto.ReadAccountsRequest{
			AccountUsername: &username,
		}
		_, err := accountQueryRepo.ReadFirst(requestDto)
		if err != nil {
			t.Errorf(
				"Expecting no error for %s, but got %s", username.String(), err.Error(),
			)
		}
	})

	t.Run("ReadFirstInvalid", func(t *testing.T) {
		username, _ := valueObject.NewUsername("invalid")
		requestDto := dto.ReadAccountsRequest{
			AccountUsername: &username,
		}

		_, err := accountQueryRepo.ReadFirst(requestDto)
		if err == nil {
			t.Errorf("Expecting error for %s, but got nil", username.String())
		}
	})

	deleteDummyUser()
}
