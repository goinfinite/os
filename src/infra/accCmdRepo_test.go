package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func TestAccCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("AddValidAccount", func(t *testing.T) {
		username := valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME"))
		password := valueObject.NewPasswordPanic(os.Getenv("DUMMY_USER_PASS"))

		addUser := dto.AddUser{
			Username: username,
			Password: password,
		}

		accCmdRepo := AccCmdRepo{}
		err := accCmdRepo.Add(addUser)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("AddInvalidAccount", func(t *testing.T) {
		username := valueObject.NewUsernamePanic("root")
		password := valueObject.NewPasswordPanic("invalid")

		addUser := dto.AddUser{
			Username: username,
			Password: password,
		}

		accCmdRepo := AccCmdRepo{}
		err := accCmdRepo.Add(addUser)
		if err == nil {
			t.Error("ExpectingError")
		}
	})
}
