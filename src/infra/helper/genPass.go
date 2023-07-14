package infraHelper

import (
	"crypto/rand"
	"math/big"
)

func GenPass(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLen := big.NewInt(int64(len(charset)))

	pass := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, charsetLen)
		pass[i] = charset[randomIndex.Int64()]
	}

	return string(pass)
}
