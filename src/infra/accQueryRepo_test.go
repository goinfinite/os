package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func TestAccQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("GetValidAccount", func(t *testing.T) {
		username := valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME"))

		authQueryRepo := AccQueryRepo{}
		accDetails, err := authQueryRepo.GetAccountDetailsByUsername(username)
		if err != nil {
			t.Error("UnexpectedError")
		}

		if (accDetails.UserId.Get()) != 1000 {
			t.Error("InvalidUserId")
		}
	})

	t.Run("GetInvalidAccount", func(t *testing.T) {
		username := valueObject.NewUsernamePanic("invalid")

		authQueryRepo := AccQueryRepo{}

		_, err := authQueryRepo.GetAccountDetailsByUsername(username)
		if err == nil {
			t.Error("ExpectingError")
		}
	})
}
