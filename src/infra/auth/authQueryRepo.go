package authInfra

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

type AuthQueryRepo struct {
}

func (repo AuthQueryRepo) IsLoginValid(login dto.Login) bool {
	storedPassHash, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"getent shadow "+login.Username.String()+" | awk -F: '{print $2}'",
	)
	if err != nil {
		return false
	}

	if len(storedPassHash) == 0 {
		return false
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(storedPassHash),
		[]byte(login.Password.String()),
	)
	return err == nil
}

func (repo AuthQueryRepo) getSessionTokenClaims(
	sessionToken valueObject.AccessTokenStr,
) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(
		sessionToken.String(),
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims, areClaimsReadable := parsedToken.Claims.(jwt.MapClaims)
	if !areClaimsReadable {
		return jwt.MapClaims{}, errors.New("SessionTokenClaimsUnReadable")
	}

	return claims, nil
}

func (repo AuthQueryRepo) getTokenDetailsFromSession(
	sessionTokenClaims jwt.MapClaims,
) (dto.AccessTokenDetails, error) {
	issuedIp, err := valueObject.NewIpAddress(
		sessionTokenClaims["originalIp"].(string),
	)
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("OriginalIpUnreadable")
	}

	var accountId valueObject.AccountId
	switch id := sessionTokenClaims["accountId"].(type) {
	case string:
		accountId, err = valueObject.NewAccountIdFromString(id)
	case float64:
		accountId, err = valueObject.NewAccountIdFromFloat(id)
	}
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("AccountIdUnreadable")
	}

	return dto.NewAccessTokenDetails(
		valueObject.NewAccessTokenTypePanic("sessionToken"),
		accountId,
		&issuedIp,
	), nil
}

func (repo AuthQueryRepo) decryptApiKey(
	token valueObject.AccessTokenStr,
) (string, error) {
	apiKeyDecoded, err := base64.StdEncoding.DecodeString(
		token.String(),
	)
	if err != nil {
		return "", errors.New("ApiKeyDecodingError")
	}
	if len(apiKeyDecoded) < aes.BlockSize {
		return "", errors.New("ApiKeyTooShort")
	}

	secretKey := os.Getenv("UAK_SECRET")
	secretKeyBytes, err := base64.RawURLEncoding.DecodeString(secretKey)
	if err != nil {
		return "", errors.New("ApiKeySecretDecodingError")
	}

	block, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		return "", errors.New("ApiKeyCipherError")
	}

	apiKeyDecryptedBinary := make([]byte, len(apiKeyDecoded)-aes.BlockSize)
	iv := apiKeyDecoded[:aes.BlockSize]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(apiKeyDecryptedBinary, apiKeyDecoded[aes.BlockSize:])

	return string(apiKeyDecryptedBinary), nil
}

func (repo AuthQueryRepo) getKeyHash(
	accountId valueObject.AccountId,
) (string, error) {
	keysHashFile := ".accountApiKeys"
	if _, err := os.Stat(keysHashFile); err != nil {
		return "", errors.New("KeysHashFileUnreadable")
	}

	getKeyCmd := exec.Command(
		"sed",
		"-n",
		"/"+accountId.String()+":/p",
		keysHashFile,
	)
	getKeyOutput, err := getKeyCmd.Output()
	if err != nil {
		return "", errors.New("KeysHashFileUnreadable")
	}
	if len(getKeyOutput) == 0 {
		return "", errors.New("AccountKeyNotFound")
	}

	// lineFormat: accountId:uuidHash
	lineContent := strings.TrimSpace(string(getKeyOutput))
	lineParts := strings.Split(lineContent, ":")
	if len(lineParts) != 2 {
		return "", errors.New("AccountKeyFormatError")
	}

	return lineParts[1], nil
}

func (repo AuthQueryRepo) getTokenDetailsFromApiKey(
	token valueObject.AccessTokenStr,
) (dto.AccessTokenDetails, error) {
	decryptedApiKey, err := repo.decryptApiKey(token)
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("ApiKeyDecryptionError")
	}

	// keyFormat: accountId:UUIDv4
	keyParts := strings.Split(decryptedApiKey, ":")
	if len(keyParts) != 2 {
		return dto.AccessTokenDetails{}, errors.New("ApiKeyFormatError")
	}

	accountId, err := valueObject.NewAccountIdFromString(keyParts[0])
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("AccountIdUnreadable")
	}
	uuid := keyParts[1]

	uuidHash := sha3.New256()
	uuidHash.Write([]byte(uuid))
	uuidHashStr := hex.EncodeToString(uuidHash.Sum(nil))

	storedUuidHash, err := repo.getKeyHash(accountId)
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("AccountKeyHashUnreadable")
	}

	if uuidHashStr != storedUuidHash {
		return dto.AccessTokenDetails{}, errors.New("AccountKeyHashMismatch")
	}

	return dto.NewAccessTokenDetails(
		valueObject.NewAccessTokenTypePanic("accountApiKey"),
		accountId,
		nil,
	), nil
}

func (repo AuthQueryRepo) GetAccessTokenDetails(
	token valueObject.AccessTokenStr,
) (dto.AccessTokenDetails, error) {
	sessionTokenClaims, err := repo.getSessionTokenClaims(token)
	if err != nil {
		return repo.getTokenDetailsFromApiKey(token)
	}

	return repo.getTokenDetailsFromSession(sessionTokenClaims)
}
