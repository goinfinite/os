package authInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestAuthCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	authCmdRepo := AuthCmdRepo{}

	t.Run("GetSessionToken", func(t *testing.T) {
		token, err := authCmdRepo.CreateSessionToken(
			valueObject.AccountId(1000),
			valueObject.NewUnixTimeAfterNow(useCase.SessionTokenExpiresIn),
			valueObject.NewLocalhostIpAddress(),
		)
		if err != nil {
			t.Errorf("UnexpectedError: %s", err.Error())
		}

		if token.TokenStr == "" {
			t.Errorf("Expected token not to be empty")
		}
	})
}
