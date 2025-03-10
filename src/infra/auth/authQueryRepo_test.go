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

	accountId, _ := valueObject.NewAccountId(1001)
	username, _ := valueObject.NewUsername("authDummyUser")
	password, _ := valueObject.NewPassword("q1w2e3r4t5y6")
	localIpAddress := valueObject.NewLocalhostIpAddress()
	createDto := dto.NewCreateAccount(
		username, password, false, accountId, localIpAddress,
	)

	_, err := accountCmdRepo.Create(createDto)
	if err != nil {
		t.Fatal("FailedToCreateDummyAccount")
	}

	t.Run("ValidLoginCredentials", func(t *testing.T) {
		createDto := dto.NewCreateSessionToken(username, password, localIpAddress)
		isValid := authQueryRepo.IsLoginValid(createDto)
		if !isValid {
			t.Fatal("LoginCredentialsInvalid")
		}
	})

	t.Run("InvalidLoginCredentials", func(t *testing.T) {
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
