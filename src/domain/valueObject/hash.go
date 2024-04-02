package valueObject

import (
	"errors"
	"regexp"
)

const hashRegex string = `^\w{6,256}$`

type Hash string

func NewHash(value string) (Hash, error) {
	hash := Hash(value)
	if !hash.isValid() {
		return "", errors.New("InvalidHash")
	}

	return hash, nil
}

func NewHashPanic(value string) Hash {
	hash, err := NewHash(value)
	if err != nil {
		panic(err)
	}

	return hash
}

func (hash Hash) isValid() bool {
	re := regexp.MustCompile(hashRegex)
	return re.MatchString(string(hash))
}

func (hash Hash) String() string {
	return string(hash)
}
