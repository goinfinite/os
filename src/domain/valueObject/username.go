package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var usernameRegex = regexp.MustCompile(`^[a-z_]([a-z0-9_-]{0,31}|[a-z0-9_-]{0,30}\$)$`)

type Username string

func NewUsername(value any) (username Username, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return username, errors.New("UsernameMustBeString")
	}

	if !usernameRegex.MatchString(stringValue) {
		return username, errors.New("InvalidUsername")
	}

	return Username(stringValue), nil
}

func (vo Username) String() string {
	return string(vo)
}
