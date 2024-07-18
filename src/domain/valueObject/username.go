package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const usernameRegex string = `^[a-z_]([a-z0-9_-]{0,31}|[a-z0-9_-]{0,30}\$)$`

type Username string

func NewUsername(value interface{}) (username Username, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return username, errors.New("UsernameValueMustBeString")
	}

	re := regexp.MustCompile(usernameRegex)
	isValid := re.MatchString(stringValue)
	if !isValid {
		return "", errors.New("InvalidUsername")
	}
	return Username(stringValue), nil
}

// TODO: remove this constructor when no longer used.
func NewUsernamePanic(value interface{}) Username {
	user, err := NewUsername(value)
	if err != nil {
		panic(err)
	}
	return user
}

func (vo Username) String() string {
	return string(vo)
}
