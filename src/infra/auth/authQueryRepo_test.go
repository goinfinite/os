package authInfra

import (
	"encoding/base64"
	"os"
	"testing"
	"time"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuthQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	authQueryRepo := NewAuthQueryRepo(testHelpers.GetPersistentDbSvc())
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
			t.Error("ExpectedInvalidCredentials")
		}
	})

	t.Run("ValidSessionAccessToken", func(t *testing.T) {
		authCmdRepo := NewAuthCmdRepo()

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

	t.Run("SessionToken_CreationAndValidation_Success", func(t *testing.T) {
		authCmdRepo := NewAuthCmdRepo()
		token, err := authCmdRepo.CreateSessionToken(
			tkValueObject.AccountId(1000),
			tkValueObject.NewUnixTimeAfterNow(3*time.Hour),
			tkValueObject.IpAddressLocal,
		)
		if err != nil {
			t.Fatalf("UnexpectedCreationError: %s", err.Error())
		}

		claims, err := authQueryRepo.readSessionTokenClaims(token.TokenStr)
		if err != nil {
			t.Fatalf("UnexpectedParseError: %s", err.Error())
		}

		if claims["accountId"] == nil {
			t.Fatal("MissingAccountIdClaim")
		}
	})

	t.Run("SessionToken_ExpiredToken_ReturnsExpiredError", func(t *testing.T) {
		jwtSecret := os.Getenv("JWT_SECRET")

		pastTime := time.Now().Add(-1 * time.Hour).Unix()
		claims := jwt.MapClaims{
			"iss":        "test",
			"iat":        pastTime,
			"nbf":        pastTime,
			"exp":        pastTime,
			"accountId":  uint64(1000),
			"originalIp": "127.0.0.1",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			t.Fatalf("FailedToCreateExpiredToken: %s", err.Error())
		}

		tokenValue, _ := tkValueObject.NewAccessTokenValue(tokenStr)
		_, err = authQueryRepo.readSessionTokenClaims(tokenValue)
		if err != errSessionTokenExpired {
			t.Fatalf("ExpectedSessionTokenExpired: %v", err)
		}
	})

	t.Run("SessionToken_WrongSecret_ReturnsSignatureInvalidError", func(t *testing.T) {
		futureTime := time.Now().Add(1 * time.Hour).Unix()
		claims := jwt.MapClaims{
			"iss":        "test",
			"iat":        time.Now().Unix(),
			"nbf":        time.Now().Unix(),
			"exp":        futureTime,
			"accountId":  uint64(1000),
			"originalIp": "127.0.0.1",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte("wrong-secret-key"))
		if err != nil {
			t.Fatalf("FailedToCreateToken: %s", err.Error())
		}

		tokenValue, _ := tkValueObject.NewAccessTokenValue(tokenStr)
		_, err = authQueryRepo.readSessionTokenClaims(tokenValue)
		if err != errSessionTokenSignatureInvalid {
			t.Fatalf("ExpectedSignatureInvalid: %v", err)
		}
	})

	t.Run("SessionToken_AlgorithmSubstitutionAttack_Prevented", func(t *testing.T) {
		// Craft a token with alg "none" manually — must be rejected by WithValidMethods
		headerJSON := `{"alg":"none","typ":"JWT"}`
		payloadJSON := `{"iss":"test","iat":1,"nbf":1,"exp":9999999999,"accountId":1000,"originalIp":"127.0.0.1"}`

		encodeSegment := func(s string) string {
			return base64.RawURLEncoding.EncodeToString([]byte(s))
		}

		noneToken := encodeSegment(headerJSON) + "." + encodeSegment(payloadJSON) + "."
		tokenValue, err := tkValueObject.NewAccessTokenValue(noneToken)
		if err != nil {
			t.Fatalf("FailedToCreateTokenValue: %s", err.Error())
		}

		_, err = authQueryRepo.readSessionTokenClaims(tokenValue)
		if err == nil {
			t.Fatal("ExpectedAlgNoneRejection")
		}
		if err == errSessionTokenExpired {
			t.Fatal("AlgNoneReachedExpiryCheck")
		}
	})
}
