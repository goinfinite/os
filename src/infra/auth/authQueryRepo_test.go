package authInfra

import (
	"os"
	"testing"
	"time"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
)

func TestAuthQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	authQueryRepo := AuthQueryRepo{}
	accountCmdRepo := accountInfra.NewAccountCmdRepo(testHelpers.GetPersistentDbSvc())
	localIpAddress := valueObject.NewLocalhostIpAddress()

	t.Run("ValidLoginCredentials", func(t *testing.T) {
		username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))
		password, _ := valueObject.NewPassword(os.Getenv("DUMMY_USER_PASS"))

		createDto := dto.NewCreateSessionToken(username, password, localIpAddress)
		isValid := authQueryRepo.IsLoginValid(createDto)
		if !isValid {
			t.Fatal("LoginCredentialsInvalid")
		}
	})

	t.Run("InvalidLoginCredentials", func(t *testing.T) {
		username, _ := valueObject.NewUsername(os.Getenv("DUMMY_USER_NAME"))
		password, _ := valueObject.NewPassword("wrongPassword")

		createDto := dto.NewCreateSessionToken(username, password, localIpAddress)
		isValid := authQueryRepo.IsLoginValid(createDto)
		if isValid {
			t.Error("Expected invalid login credentials, but got valid")
		}
	})

	t.Run("ValidSessionAccessToken", func(t *testing.T) {
		authCmdRepo := AuthCmdRepo{}

		token, _ := authCmdRepo.CreateSessionToken(
			valueObject.AccountId(1000),
			valueObject.NewUnixTimeAfterNow(3*time.Hour),
			valueObject.NewLocalhostIpAddress(),
		)

		_, err := authQueryRepo.ReadAccessTokenDetails(token.TokenStr)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("InvalidSessionAccessToken", func(t *testing.T) {
		invalidToken, _ := valueObject.NewAccessTokenStr(
			"invalidTokenInvalidTokenInvalidTokenInvalidTokenInvalidToken",
		)
		_, err := authQueryRepo.ReadAccessTokenDetails(invalidToken)
		if err == nil {
			t.Error("ExpectingError")
		}
	})

	t.Run("ValidAccountApiKey", func(t *testing.T) {
		accountId, _ := valueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))
		apiKey, err := accountCmdRepo.UpdateApiKey(accountId)
		if err != nil {
			t.Error(err)
		}

		_, err = authQueryRepo.ReadAccessTokenDetails(apiKey)
		if err != nil {
			t.Error(err)
		}
	})
}
