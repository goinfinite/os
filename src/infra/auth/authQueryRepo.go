package authInfra

import (
	"crypto/sha3"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	tkInfra "github.com/goinfinite/tk/src/infra"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	errSessionTokenExpired          = errors.New("SessionTokenExpired")
	errSessionTokenSignatureInvalid = errors.New("SessionTokenSignatureInvalid")
	errSessionTokenParseError       = errors.New("SessionTokenParseError")
	errSessionTokenClaimsUnreadable = errors.New("SessionTokenClaimsUnreadable")
)

type AuthQueryRepo struct {
	persistentDbSvc         *internalDbInfra.PersistentDatabaseService
	sessionTokenSecretBytes []byte
}

func NewAuthQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *AuthQueryRepo {
	sessionTokenSecret := os.Getenv("JWT_SECRET")
	return &AuthQueryRepo{
		persistentDbSvc:         persistentDbSvc,
		sessionTokenSecretBytes: []byte(sessionTokenSecret),
	}
}

func (repo *AuthQueryRepo) IsLoginValid(createDto dto.CreateSessionToken) bool {
	readStoredPassHashCmd := "getent shadow " + createDto.Username.String() +
		" | awk -F: '{print $2}'"
	storedPassHash, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:            readStoredPassHashCmd,
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		slog.Debug(
			"GetentShadowError",
			slog.String("username", createDto.Username.String()),
			slog.String("err", err.Error()),
		)
		return false
	}

	if len(storedPassHash) == 0 {
		return false
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(storedPassHash),
		[]byte(createDto.Password.String()),
	)
	return err == nil
}

func (repo *AuthQueryRepo) readSessionTokenClaims(
	sessionToken tkValueObject.AccessTokenValue,
) (claims jwt.MapClaims, err error) {
	parsedToken, err := jwt.Parse(
		sessionToken.String(),
		func(token *jwt.Token) (interface{}, error) {
			return repo.sessionTokenSecretBytes, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return claims, errSessionTokenExpired
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return claims, errSessionTokenSignatureInvalid
		default:
			return claims, errSessionTokenParseError
		}
	}

	if !parsedToken.Valid {
		return claims, errSessionTokenParseError
	}

	claims, areClaimsReadable := parsedToken.Claims.(jwt.MapClaims)
	if !areClaimsReadable {
		return claims, errSessionTokenClaimsUnreadable
	}

	return claims, nil
}

func (repo *AuthQueryRepo) readTokenDetailsFromSession(
	sessionTokenClaims jwt.MapClaims,
) (tokenDetails dto.AccessTokenDetails, err error) {
	tokenType, _ := tkValueObject.NewAccessTokenType("sessionToken")

	accountId, err := tkValueObject.NewAccountId(sessionTokenClaims["accountId"])
	if err != nil {
		return tokenDetails, errors.New("AccountIdUnreadable")
	}

	issuedIp, err := tkValueObject.NewIpAddress(sessionTokenClaims["originalIp"])
	if err != nil {
		return tokenDetails, errors.New("OriginalIpUnreadable")
	}

	return dto.NewAccessTokenDetails(tokenType, accountId, &issuedIp), nil
}

func (repo *AuthQueryRepo) readKeyHash(
	accountId tkValueObject.AccountId,
) (keyHash string, err error) {
	accountModel := dbModel.Account{ID: accountId.Uint64()}
	err = repo.persistentDbSvc.Handler.Model(&accountModel).First(&accountModel).Error
	if err != nil {
		return keyHash, errors.New("AccountNotFound")
	}

	if accountModel.KeyHash == nil {
		return keyHash, errors.New("UserKeyHashNotFound")
	}

	return *accountModel.KeyHash, nil
}

func (repo *AuthQueryRepo) readTokenDetailsFromApiKey(
	token tkValueObject.AccessTokenValue,
) (tokenDetails dto.AccessTokenDetails, err error) {
	secretKey := os.Getenv("ACCOUNT_API_KEY_SECRET")
	cypher, err := tkInfra.NewCypher(secretKey)
	if err != nil {
		return tokenDetails, errors.New("ApiKeyDecryptSecretKeyError")
	}
	decryptedApiKey, err := cypher.Decrypt(token.String())
	if err != nil {
		return tokenDetails, errors.New("ApiKeyDecryptionError")
	}

	// keyFormat: accountId:UUIDv4
	keyParts := strings.Split(decryptedApiKey, ":")
	if len(keyParts) != 2 {
		return tokenDetails, errors.New("ApiKeyFormatError")
	}

	accountId, err := tkValueObject.NewAccountId(keyParts[0])
	if err != nil {
		return tokenDetails, errors.New("AccountIdUnreadable")
	}

	apiKeyHasher := sha3.New256()
	apiKeyHasher.Write([]byte(decryptedApiKey))
	apiKeyHashStr := hex.EncodeToString(apiKeyHasher.Sum(nil))

	storedUuidHash, err := repo.readKeyHash(accountId)
	if err != nil {
		return tokenDetails, errors.New("UserKeyHashUnreadable")
	}

	if subtle.ConstantTimeCompare([]byte(apiKeyHashStr), []byte(storedUuidHash)) != 1 {
		return tokenDetails, errors.New("UserKeyHashMismatch")
	}

	tokenType, _ := tkValueObject.NewAccessTokenType("accountApiKey")

	return dto.NewAccessTokenDetails(tokenType, accountId, nil), nil
}

func (repo *AuthQueryRepo) ReadAccessTokenDetails(
	token tkValueObject.AccessTokenValue,
) (tokenDetails dto.AccessTokenDetails, err error) {
	sessionTokenClaims, err := repo.readSessionTokenClaims(token)
	if err != nil {
		isLikelyApiKey := err == errSessionTokenParseError
		if !isLikelyApiKey {
			return tokenDetails, err
		}

		return repo.readTokenDetailsFromApiKey(token)
	}

	return repo.readTokenDetailsFromSession(sessionTokenClaims)
}
