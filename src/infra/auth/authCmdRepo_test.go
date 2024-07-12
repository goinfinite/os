package authInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestAuthCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("GetSessionToken", func(t *testing.T) {
		authCmdRepo := AuthCmdRepo{}
		token, err := authCmdRepo.GenerateSessionToken(
			valueObject.AccountId(1000),
			valueObject.NewUnixTimeAfterNow(useCase.SessionTokenExpiresIn),
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
