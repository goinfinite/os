package infraHelper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
)

func EncryptStr(
	secretKey, plainText string,
) (encryptedText string, err error) {
	secretKeyBytes, err := base64.RawURLEncoding.DecodeString(secretKey)
	if err != nil {
		slog.Error("EncryptSecretKeyError", slog.Any("error", err))
		return encryptedText, errors.New("EncryptSecretKeyError")
	}

	block, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		slog.Error("EncryptCipherError", slog.Any("error", err))
		return encryptedText, errors.New("EncryptCipherError")
	}

	plainTextBytes := []byte(plainText)
	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		slog.Error("EncryptIvGenerationError", slog.Any("error", err))
		return encryptedText, errors.New("EncryptIvGenerationError")
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptStr(
	secretKey, encryptedText string,
) (decryptedText string, err error) {
	apiKeyDecoded, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		slog.Error("DecryptDecodingError", slog.Any("error", err))
		return decryptedText, errors.New("DecryptDecodingError")
	}
	if len(apiKeyDecoded) < aes.BlockSize {
		return decryptedText, errors.New("DecryptDecodedTooShort")
	}

	secretKeyBytes, err := base64.RawURLEncoding.DecodeString(secretKey)
	if err != nil {
		slog.Error("DecryptSecretDecodingError", slog.Any("error", err))
		return decryptedText, errors.New("DecryptSecretDecodingError")
	}

	block, err := aes.NewCipher(secretKeyBytes)
	if err != nil {
		slog.Error("DecryptCipherError", slog.Any("error", err))
		return decryptedText, errors.New("DecryptCipherError")
	}

	apiKeyDecryptedBinary := make([]byte, len(apiKeyDecoded)-aes.BlockSize)
	iv := apiKeyDecoded[:aes.BlockSize]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(apiKeyDecryptedBinary, apiKeyDecoded[aes.BlockSize:])

	return string(apiKeyDecryptedBinary), nil
}
