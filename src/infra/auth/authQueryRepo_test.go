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

	accountId, _ := valueObject.NewAccountId(1000)
	expiresIn := valueObject.NewUnixTimeAfterNow(useCase.SessionTokenExpiresIn)
	localIpAddress := valueObject.NewLocalhostIpAddress()
	token, err := authCmdRepo.GenerateSessionToken(
		accountId, expiresIn, localIpAddress,
	)
	if err != nil {
		t.Errorf("UnexpectedError: %s", err.Error())
	}

	username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))

	t.Run("ValidLoginCredentials", func(t *testing.T) {
		password, _ := valueObject.NewPassword(os.Getenv("DUMMY_USER_PASS"))
		login := dto.NewLogin(username, password, localIpAddress)

		isValid := authQueryRepo.IsLoginValid(login)
		if !isValid {
			t.Error("Expected valid login credentials, but got invalid")
		}
	})

	t.Run("InvalidLoginCredentials", func(t *testing.T) {
		password, _ := valueObject.NewPassword("wrongPassword")
		login := dto.NewLogin(username, password, localIpAddress)

		isValid := authQueryRepo.IsLoginValid(login)
		if isValid {
			t.Error("Expected invalid login credentials, but got valid")
		}
	})

	t.Run("ValidSessionAccessToken", func(t *testing.T) {
		_, err = authQueryRepo.GetAccessTokenDetails(token.TokenStr)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("InvalidSessionAccessToken", func(t *testing.T) {
		invalidToken, _ := valueObject.NewAccessTokenStr(
			"invalidTokenInvalidTokenInvalidTokenInvalidTokenInvalidToken",
		)
		_, err := authQueryRepo.GetAccessTokenDetails(invalidToken)
		if err == nil {
			t.Error("ExpectingError")
		}
	})

	t.Run("DecryptValidApiKey", func(t *testing.T) {
		tokenBytes := []byte(token.TokenStr.String())
		apiKeyStr := base64.StdEncoding.EncodeToString(tokenBytes)
		apiKey, _ := valueObject.NewAccessTokenStr(apiKeyStr)

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
