package infraHelper

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log/slog"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

func EncryptStr(
	secretKey, plainText string,
) (encryptedText string, err error) {
	cypher, err := tkInfra.NewCypher(secretKey)
	if err != nil {
		return encryptedText, errors.New("EncryptSecretKeyError")
	}

	return cypher.Encrypt(plainText)
}

func legacyCtrDecrypt(
	secretKey, encryptedText string,
) (decryptedText string, err error) {
	apiKeyDecoded, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return decryptedText, errors.New("DecryptDecodingError")
	}
	if len(apiKeyDecoded) < aes.BlockSize {
		return decryptedText, errors.New("DecryptDecodedTooShort")
	}

	secretKeyBytes, err := base64.RawURLEncoding.DecodeString(secretKey)
	if err != nil {
		return decryptedText, errors.New("DecryptSecretDecodingError")
	}

	block, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		return decryptedText, errors.New("DecryptCipherError")
	}

	apiKeyDecryptedBinary := make([]byte, len(apiKeyDecoded)-aes.BlockSize)
	iv := apiKeyDecoded[:aes.BlockSize]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(apiKeyDecryptedBinary, apiKeyDecoded[aes.BlockSize:])

	return string(apiKeyDecryptedBinary), nil
}

func DecryptStr(
	secretKey, encryptedText string,
) (decryptedText string, err error) {
	cypher, err := tkInfra.NewCypher(secretKey)
	if err != nil {
		return decryptedText, errors.New("DecryptSecretKeyError")
	}

	decryptedText, err = cypher.Decrypt(encryptedText)
	if err == nil {
		return decryptedText, nil
	}

	slog.Debug("GcmDecryptFailed, trying legacy CTR decryption")
	return legacyCtrDecrypt(secretKey, encryptedText)
}
