package authInfra

import (
	"testing"
	"time"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestAuthCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("GetSessionToken", func(t *testing.T) {
		authCmdRepo := AuthCmdRepo{}
		token, err := authCmdRepo.GenerateSessionToken(
			valueObject.AccountId(1000),
			valueObject.UnixTime(
				time.Now().Add(3*time.Hour).Unix(),
			),
			valueObject.NewIpAddressPanic("127.0.0.1"),
		)
		if err != nil {
			t.Errorf("UnexpectedError: %s", err.Error())
		}

		if token.TokenStr == "" {
			t.Errorf("Expected token not to be empty")
		}
	})
}
