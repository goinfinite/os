package authInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/useCase"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestAuthCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	authCmdRepo := NewAuthCmdRepo()

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
			t.Errorf("EmptyToken")
		}
	})

	t.Run("SessionToken_CreationAndValidation_Success", func(t *testing.T) {
		token, err := authCmdRepo.CreateSessionToken(
			tkValueObject.AccountId(1000),
			tkValueObject.NewUnixTimeAfterNow(useCase.SessionTokenExpiresIn),
			tkValueObject.IpAddressLocal,
		)
		if err != nil {
			t.Fatalf("UnexpectedError: %s", err.Error())
		}

		if token.TokenStr.String() == "" {
			t.Fatal("EmptyTokenString")
		}

		authQueryRepo := NewAuthQueryRepo(nil)
		_, err = authQueryRepo.readSessionTokenClaims(token.TokenStr)
		if err != nil {
			t.Fatalf("TokenParseFailed: %s", err.Error())
		}
	})
}
