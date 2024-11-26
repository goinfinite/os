package voHelper

import (
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/sha3"
)

func TransformPlainContentIntoStrongHash(
	plainContent string,
) (strongHash string, err error) {
	hash := sha3.New256()
	_, err = hash.Write([]byte(plainContent))
	if err != nil {
		return strongHash, errors.New("InvalidPlainContentToHash")
	}
	encodedContentBytes := hash.Sum(nil)
	return hex.EncodeToString(encodedContentBytes), nil
}
