package infraHelper

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

func GenWeakHash(value string) string {
	hash := md5.Sum([]byte(value))
	return hex.EncodeToString(hash[:])
}

func GenStrongHash(value string) string {
	hash := sha3.New256()
	hash.Write([]byte(value))
	return hex.EncodeToString(hash.Sum(nil))
}

func GenStrongShortHash(value string) string {
	return GenStrongHash(string(value))[:12]
}
