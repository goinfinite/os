package valueObject

import (
	"errors"
	"regexp"
)

const hashRegex string = `^\w{6,256}$`

type Hash string

func NewHash(value string) (Hash, error) {
	user := Hash(value)
	if !user.isValid() {
		return "", errors.New("InvalidHash")
	}

	return user, nil
}

func NewHashPanic(value string) Hash {
	user, err := NewHash(value)
	if err != nil {
		panic(err)
	}

	return user
}

func (user Hash) isValid() bool {
	re := regexp.MustCompile(hashRegex)
	return re.MatchString(string(user))
}

func (user Hash) String() string {
	return string(user)
}
