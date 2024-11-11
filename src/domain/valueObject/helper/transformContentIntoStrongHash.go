package voHelper

import (
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/sha3"
)

func TransformContentIntoStrongHash(content string) (strongHash string, err error) {
	hash := sha3.New256()
	_, err = hash.Write([]byte(content))
	if err != nil {
		return strongHash, errors.New("InvalidContentToHash")
	}
	encodedContentBytes := hash.Sum(nil)
	return hex.EncodeToString(encodedContentBytes), nil
}
