package infra

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/msteinert/pam"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
	"golang.org/x/crypto/sha3"
)

type AuthQueryRepo struct {
}

func (repo AuthQueryRepo) IsLoginValid(login dto.Login) bool {
	tx, err := pam.StartFunc(
		"system-auth",
		login.Username.String(),
		func(s pam.Style, msg string) (string, error) {
			switch s {
			case pam.PromptEchoOff:
				return login.Password.String(), nil
			case pam.PromptEchoOn, pam.ErrorMsg, pam.TextInfo:
				log.Printf("PamMessage: %s", msg)
				return "", nil
			}
			log.Println("UnhandledPamMessageStyle")
			return "", errors.New("UnhandledPamMessageStyle")
		})

	if err != nil {
		log.Printf("PamError: %v", err)
		return false
	}

	if err = tx.Authenticate(0); err != nil {
		return false
	}

	return true
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

	var userId valueObject.UserId
	switch id := sessionTokenClaims["userId"].(type) {
	case string:
		userId, err = valueObject.NewUserIdFromString(id)
	case float64:
		userId, err = valueObject.NewUserIdFromFloat(id)
	}
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("UserIdUnreadable")
	}

	return dto.NewAccessTokenDetails(
		valueObject.NewAccessTokenTypePanic("sessionToken"),
		userId,
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
	userId valueObject.UserId,
) (string, error) {
	keysHashFile := ".userKeys"
	if _, err := os.Stat(keysHashFile); err != nil {
		return "", errors.New("KeysHashFileUnreadable")
	}

	getKeyCmd := exec.Command(
		"sed",
		"-n",
		"/"+userId.String()+":/p",
		keysHashFile,
	)
	getKeyOutput, err := getKeyCmd.Output()
	if err != nil {
		return "", errors.New("KeysHashFileUnreadable")
	}
	if len(getKeyOutput) == 0 {
		return "", errors.New("UserKeyNotFound")
	}

	// lineFormat: userId:uuidHash
	lineContent := strings.TrimSpace(string(getKeyOutput))
	lineParts := strings.Split(lineContent, ":")
	if len(lineParts) != 2 {
		return "", errors.New("UserKeyFormatError")
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

	// keyFormat: userId:UUIDv4
	keyParts := strings.Split(decryptedApiKey, ":")
	if len(keyParts) != 2 {
		return dto.AccessTokenDetails{}, errors.New("ApiKeyFormatError")
	}

	userId, err := valueObject.NewUserIdFromString(keyParts[0])
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("UserIdUnreadable")
	}
	uuid := keyParts[1]

	uuidHash := sha3.New256()
	uuidHash.Write([]byte(uuid))
	uuidHashStr := hex.EncodeToString(uuidHash.Sum(nil))

	storedUuidHash, err := repo.getKeyHash(userId)
	if err != nil {
		return dto.AccessTokenDetails{}, errors.New("UserKeyHashUnreadable")
	}

	if uuidHashStr != storedUuidHash {
		return dto.AccessTokenDetails{}, errors.New("UserKeyHashMismatch")
	}

	return dto.NewAccessTokenDetails(
		valueObject.NewAccessTokenTypePanic("userApiKey"),
		userId,
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
