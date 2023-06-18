package valueObject

import (
	"errors"
	"regexp"
)

const usernameRegex string = `^[a-z_]([a-z0-9_-]{0,31}|[a-z0-9_-]{0,30}\$)$`

type Username string

func NewUsername(value string) (Username, error) {
	user := Username(value)
	if !user.isValid() {
		return "", errors.New("InvalidUsername")
	}
	return user, nil
}

func NewUsernamePanic(value string) Username {
	user := Username(value)
	if !user.isValid() {
		panic("InvalidUsername")
	}
	return user
}

func (user Username) isValid() bool {
	re := regexp.MustCompile(usernameRegex)
	return re.MatchString(string(user))
}

func (user Username) String() string {
	return string(user)
}
