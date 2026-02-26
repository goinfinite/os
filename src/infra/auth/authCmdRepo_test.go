package authInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/useCase"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestAuthCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	authCmdRepo := AuthCmdRepo{}

	t.Run("GetSessionToken", func(t *testing.T) {
		token, err := authCmdRepo.CreateSessionToken(
			tkValueObject.AccountId(1000),
			tkValueObject.NewUnixTimeAfterNow(useCase.SessionTokenExpiresIn),
			tkValueObject.IpAddressLocal,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %s", err.Error())
		}

		if token.TokenStr == "" {
			t.Errorf("Expected token not to be empty")
		}
	})
}
