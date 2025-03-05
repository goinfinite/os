package authInfra

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewAuthQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *AuthQueryRepo {
	return &AuthQueryRepo{
		persistentDbSvc: persistentDbSvc,
	}
}

func (repo *AuthQueryRepo) IsLoginValid(createDto dto.CreateSessionToken) bool {
	readStoredPassHashCmd := "getent shadow " + createDto.Username.String() +
		" | awk -F: '{print $2}'"
	storedPassHash, err := infraHelper.RunCmd(infraHelper.RunCmdConfigs{
		Command:               readStoredPassHashCmd,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		slog.Debug(
			"GetentShadowError",
			slog.String("username", createDto.Username.String()),
			slog.Any("error", err),
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
	sessionToken valueObject.AccessTokenStr,
) (claims jwt.MapClaims, err error) {
	parsedToken, err := jwt.Parse(
		sessionToken.String(),
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
	if err != nil {
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return claims, errors.New("SessionTokenExpired")
		}

		return claims, errors.New("SessionTokenParseError: " + err.Error())
	}

	claims, areClaimsReadable := parsedToken.Claims.(jwt.MapClaims)
	if !areClaimsReadable {
		return claims, errors.New("SessionTokenClaimsUnreadable")
	}

	return claims, nil
}

func (repo *AuthQueryRepo) readTokenDetailsFromSession(
	sessionTokenClaims jwt.MapClaims,
) (tokenDetails dto.AccessTokenDetails, err error) {
	tokenType, _ := valueObject.NewAccessTokenType("sessionToken")

	accountId, err := valueObject.NewAccountId(sessionTokenClaims["accountId"])
	if err != nil {
		return tokenDetails, errors.New("AccountIdUnreadable")
	}

	issuedIp, err := valueObject.NewIpAddress(sessionTokenClaims["originalIp"])
	if err != nil {
		return tokenDetails, errors.New("OriginalIpUnreadable")
	}

	return dto.NewAccessTokenDetails(tokenType, accountId, &issuedIp), nil
}

func (repo *AuthQueryRepo) readKeyHash(
	accountId valueObject.AccountId,
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
	token valueObject.AccessTokenStr,
) (tokenDetails dto.AccessTokenDetails, err error) {
	secretKey := os.Getenv("ACCOUNT_API_KEY_SECRET")
	decryptedApiKey, err := infraHelper.DecryptStr(secretKey, token.String())
	if err != nil {
		return tokenDetails, errors.New("ApiKeyDecryptionError")
	}

	// keyFormat: accountId:UUIDv4
	keyParts := strings.Split(decryptedApiKey, ":")
	if len(keyParts) != 2 {
		return tokenDetails, errors.New("ApiKeyFormatError")
	}

	accountId, err := valueObject.NewAccountId(keyParts[0])
	if err != nil {
		return tokenDetails, errors.New("AccountIdUnreadable")
	}

	uuidHash := infraHelper.GenStrongHash(keyParts[1])

	storedUuidHash, err := repo.readKeyHash(accountId)
	if err != nil {
		return tokenDetails, errors.New("UserKeyHashUnreadable")
	}

	if uuidHash != storedUuidHash {
		return tokenDetails, errors.New("UserKeyHashMismatch")
	}

	tokenType, _ := valueObject.NewAccessTokenType("accountApiKey")

	return dto.NewAccessTokenDetails(tokenType, accountId, nil), nil
}

func (repo *AuthQueryRepo) ReadAccessTokenDetails(
	token valueObject.AccessTokenStr,
) (tokenDetails dto.AccessTokenDetails, err error) {
	sessionTokenClaims, err := repo.readSessionTokenClaims(token)
	if err != nil {
		if err.Error() == "SessionTokenExpired" {
			return tokenDetails, err
		}

		return repo.readTokenDetailsFromApiKey(token)
	}

	return repo.readTokenDetailsFromSession(sessionTokenClaims)
}
