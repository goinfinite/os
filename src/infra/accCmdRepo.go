package infra

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/google/uuid"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

type AccCmdRepo struct {
}

func (repo AccCmdRepo) Add(addUser dto.AddUser) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(addUser.Password.String()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("PasswordHashError: %s", err)
		return errors.New("PasswordHashError")
	}

	addUserCmd := exec.Command(
		"useradd",
		"-m",
		"-s", "/bin/bash",
		"-p", string(passHash),
		addUser.Username.String(),
	)

	err = addUserCmd.Run()
	if err != nil {
		log.Printf("UserAddError: %s", err)
		return errors.New("UserAddError")
	}

	return nil
}

func getUsernameById(userId valueObject.UserId) (valueObject.Username, error) {
	accQuery := AccQueryRepo{}
	accDetails, err := accQuery.GetById(userId)
	if err != nil {
		log.Printf("GetUserDetailsError: %s", err)
		return "", errors.New("GetUserDetailsError")
	}

	return accDetails.Username, nil
}

func (repo AccCmdRepo) Delete(userId valueObject.UserId) error {
	username, err := getUsernameById(userId)
	if err != nil {
		return err
	}

	delUserCmd := exec.Command(
		"userdel",
		"-r",
		username.String(),
	)

	err = delUserCmd.Run()
	if err != nil {
		log.Printf("UserDeleteError: %s", err)
		return errors.New("UserDeleteError")
	}

	return nil
}

func (repo AccCmdRepo) UpdatePassword(
	userId valueObject.UserId,
	password valueObject.Password,
) error {
	passHash, err := bcrypt.GenerateFromPassword(
		[]byte(password.String()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("PasswordHashError: %s", err)
		return errors.New("PasswordHashError")
	}

	username, err := getUsernameById(userId)
	if err != nil {
		return err
	}

	updateUserCmd := exec.Command(
		"usermod",
		"-p", string(passHash),
		username.String(),
	)

	err = updateUserCmd.Run()
	if err != nil {
		log.Printf("PasswordUpdateError: %s", err)
		return errors.New("PasswordUpdateError")
	}

	return nil
}

func encryptApiKey(plainTextApiKey string) (valueObject.AccessTokenStr, error) {
	secretKey := os.Getenv("UAK_SECRET")
	secretKeyBytes, err := base64.RawURLEncoding.DecodeString(secretKey)
	if err != nil {
		log.Printf("ApiKeySecretKeyError: %s", err)
		return "", errors.New("ApiKeySecretKeyError")
	}

	block, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		log.Printf("ApiKeyCipherError: %s", err)
		return "", errors.New("ApiKeyCipherError")
	}

	plainTextApiKeyBytes := []byte(plainTextApiKey)
	cipherText := make([]byte, aes.BlockSize+len(plainTextApiKeyBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Printf("ApiKeyIvGenerationError: %s", err)
		return "", errors.New("ApiKeyIvGenerationError")
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextApiKeyBytes)

	newApiKey, err := valueObject.NewAccessTokenStr(
		base64.StdEncoding.EncodeToString(cipherText),
	)
	if err != nil {
		log.Printf("ApiKeyEncodingError: %s", err)
		return "", errors.New("ApiKeyEncodingError")
	}

	return newApiKey, nil
}

func storeNewApiKeyHash(
	userId valueObject.UserId,
	apiKey valueObject.AccessTokenStr,
) error {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()

	hashKeysFile := ".userKeys"

	if _, err := os.Stat(hashKeysFile); err == nil {
		removeOldKeyCmd := exec.Command(
			"sed",
			"-i",
			"/"+userId.String()+"/d",
			hashKeysFile,
		)
		err := removeOldKeyCmd.Run()
		if err != nil {
			log.Printf("UserKeysRemoveOldKeyError: %s", err)
			return errors.New("UserKeysRemoveOldKeyError")
		}
	}

	hash := sha3.New256()
	hash.Write([]byte(apiKey.String()))
	hashString := hex.EncodeToString(hash.Sum(nil))

	file, err := os.OpenFile(hashKeysFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0400)
	if err != nil {
		log.Printf("UserKeysOpenFileError: %v", err)
		return errors.New("UserKeysOpenFileError")
	}
	defer file.Close()

	_, err = file.WriteString(userId.String() + ":" + hashString + "\n")
	if err != nil {
		log.Printf("UserKeysWriteError: %v", err)
		return errors.New("UserKeysWriteError")
	}

	return nil
}

func (repo AccCmdRepo) UpdateApiKey(
	userId valueObject.UserId,
) (valueObject.AccessTokenStr, error) {
	plainTextApiKey := userId.String() + ":" + uuid.New().String()
	newApiKey, err := encryptApiKey(plainTextApiKey)
	if err != nil {
		return "", err
	}

	err = storeNewApiKeyHash(userId, newApiKey)
	if err != nil {
		return "", err
	}

	return newApiKey, nil
}
