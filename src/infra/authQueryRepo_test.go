package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func TestAuthQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("ValidLoginCredentials", func(t *testing.T) {
		login := dto.NewLogin(
			valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME")),
			valueObject.NewPasswordPanic(os.Getenv("DUMMY_USER_PASS")),
		)
		authQueryRepo := AuthQueryRepo{}
		isValid := authQueryRepo.IsLoginValid(login)
		if !isValid {
			t.Error("Expected valid login credentials, but got invalid")
		}
	})

	t.Run("InvalidLoginCredentials", func(t *testing.T) {
		login := dto.NewLogin(
			valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME")),
			valueObject.NewPasswordPanic("wrongPassword"),
		)
		authQueryRepo := AuthQueryRepo{}
		isValid := authQueryRepo.IsLoginValid(login)
		if isValid {
			t.Error("Expected invalid login credentials, but got valid")
		}
	})
}
