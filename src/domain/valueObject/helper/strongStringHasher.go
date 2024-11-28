package voHelper

import (
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/sha3"
)

func StrongStringHasher(
	stringToHash string,
) (strongHash string, err error) {
	hash := sha3.New256()
	_, err = hash.Write([]byte(stringToHash))
	if err != nil {
		return strongHash, errors.New("InvalidStringToHash")
	}
	encodedContentBytes := hash.Sum(nil)
	return hex.EncodeToString(encodedContentBytes), nil
}
