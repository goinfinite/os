package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

const usernameRegex string = `^[a-z_]([a-z0-9_-]{0,31}|[a-z0-9_-]{0,30}\$)$`

type Username string

func NewUsername(value interface{}) (username Username, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return username, errors.New("UsernameMustBeString")
	}

	re := regexp.MustCompile(usernameRegex)
	if !re.MatchString(stringValue) {
		return username, errors.New("InvalidUsername")
	}

	return Username(stringValue), nil
}

func (vo Username) String() string {
	return string(vo)
}
