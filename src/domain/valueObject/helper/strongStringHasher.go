package voHelper

import (
	"crypto/sha3"
	"encoding/hex"
)

const shortHashLength = 12

func StrongStringHasher(stringToHash string) (strongHash string) {
	digest := sha3.Sum256([]byte(stringToHash))
	return hex.EncodeToString(digest[:])
}

func StrongStringShortHasher(stringToHash string) (shortHash string) {
	return StrongStringHasher(stringToHash)[:shortHashLength]
}
