package authInfra

import (
	"os"
	"testing"
	"time"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestAuthQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	authQueryRepo := AuthQueryRepo{}
	accountCmdRepo := accountInfra.NewAccountCmdRepo(testHelpers.GetPersistentDbSvc())

	accountId, _ := tkValueObject.NewAccountId(1001)
	username, _ := valueObject.NewUsername("authDummyUser")
	rawPassword := "q1w2e3r4!5y6"
	accountPassword, _ := tkValueObject.NewPassword(rawPassword)
	sessionPassword, _ := tkValueObject.NewWeakPassword(rawPassword)
	localIpAddress := tkValueObject.IpAddressLocal
	createDto := dto.NewCreateAccount(
		username, accountPassword, false, accountId, localIpAddress,
	)

	_, err := accountCmdRepo.Create(createDto)
	if err != nil {
		t.Fatal("FailedToCreateDummyAccount")
	}

	t.Run("ValidLoginCredentials", func(t *testing.T) {
		createDto := dto.NewCreateSessionToken(username, sessionPassword, localIpAddress)
		isValid := authQueryRepo.IsLoginValid(createDto)
		if !isValid {
			t.Fatal("LoginCredentialsInvalid")
		}
	})

	t.Run("InvalidLoginCredentials", func(t *testing.T) {
		wrongPassword, _ := tkValueObject.NewWeakPassword("wrongPassword")

		createDto := dto.NewCreateSessionToken(username, wrongPassword, localIpAddress)
		isValid := authQueryRepo.IsLoginValid(createDto)
		if isValid {
			t.Error("Expected invalid login credentials, but got valid")
		}
	})

	t.Run("ValidSessionAccessToken", func(t *testing.T) {
		authCmdRepo := AuthCmdRepo{}

		token, _ := authCmdRepo.CreateSessionToken(
			tkValueObject.AccountId(1000),
			tkValueObject.NewUnixTimeAfterNow(3*time.Hour),
			tkValueObject.IpAddressLocal,
		)

		_, err := authQueryRepo.ReadAccessTokenDetails(token.TokenStr)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("InvalidSessionAccessToken", func(t *testing.T) {
		invalidToken, _ := tkValueObject.NewAccessTokenValue(
			"invalidTokenInvalidTokenInvalidTokenInvalidTokenInvalidToken",
		)
		_, err := authQueryRepo.ReadAccessTokenDetails(invalidToken)
		if err == nil {
			t.Error("ExpectingError")
		}
	})

	t.Run("ValidAccountApiKey", func(t *testing.T) {
		accountId, _ := tkValueObject.NewAccountId(os.Getenv("DUMMY_USER_ID"))
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
