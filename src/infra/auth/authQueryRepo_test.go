package authInfra

import (
	"encoding/base64"
	"os"
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestAuthQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	authQueryRepo := AuthQueryRepo{}
	authCmdRepo := AuthCmdRepo{}

	token, err := authCmdRepo.GenerateSessionToken(
		valueObject.AccountId(1000),
		valueObject.NewUnixTimeAfterNow(useCase.SessionTokenExpiresIn),
		valueObject.NewLocalhostIpAddress(),
	)
	if err != nil {
		t.Errorf("UnexpectedError: %s", err.Error())
	}

	t.Run("ValidLoginCredentials", func(t *testing.T) {
		login := dto.NewLogin(
			valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME")),
			valueObject.NewPasswordPanic(os.Getenv("DUMMY_USER_PASS")),
			valueObject.NewLocalhostIpAddress(),
		)

		isValid := authQueryRepo.IsLoginValid(login)
		if !isValid {
			t.Error("Expected valid login credentials, but got invalid")
		}
	})

	t.Run("InvalidLoginCredentials", func(t *testing.T) {
		login := dto.NewLogin(
			valueObject.NewUsernamePanic(os.Getenv("DUMMY_USER_NAME")),
			valueObject.NewPasswordPanic("wrongPassword"),
			valueObject.NewLocalhostIpAddress(),
		)

		isValid := authQueryRepo.IsLoginValid(login)
		if isValid {
			t.Error("Expected invalid login credentials, but got valid")
		}
	})

	t.Run("ValidSessionAccessToken", func(t *testing.T) {
		_, err = authQueryRepo.ReadAccessTokenDetails(token.TokenStr)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("InvalidSessionAccessToken", func(t *testing.T) {
		invalidToken := valueObject.NewAccessTokenStrPanic(
			"invalidTokenInvalidTokenInvalidTokenInvalidTokenInvalidToken",
		)
		_, err := authQueryRepo.ReadAccessTokenDetails(invalidToken)
		if err == nil {
			t.Error("ExpectingError")
		}
	})

	t.Run("DecryptValidApiKey", func(t *testing.T) {
		tokenBytes := []byte(token.TokenStr.String())
		apiKeyStr := base64.StdEncoding.EncodeToString(tokenBytes)
		apiKey := valueObject.NewAccessTokenStrPanic(apiKeyStr)

		_, err := authQueryRepo.decryptApiKey(apiKey)
		if err != nil {
			t.Errorf(
				"Unexpected '%s' error for '%s'",
				err.Error(),
				apiKeyStr,
			)
		}
	})

	t.Run("DecryptInvalidApiKey", func(t *testing.T) {
		_, err := authQueryRepo.decryptApiKey(token.TokenStr)
		if err == nil {
			t.Errorf("Expecting error for '%s'", token.TokenStr.String())
		}
	})
}
