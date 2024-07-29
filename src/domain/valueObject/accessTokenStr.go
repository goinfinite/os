package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const accessTokenStrRegex = `^[a-zA-Z0-9\-_=+/.]{22,444}$`

type AccessTokenStr string

func NewAccessTokenStr(value interface{}) (accessTokenStr AccessTokenStr, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return accessTokenStr, errors.New("AccessTokenStrMustBeString")
	}

	re := regexp.MustCompile(accessTokenStrRegex)
	if !re.MatchString(stringValue) {
		return "", errors.New("InvalidAccessTokenStr")
	}

	return AccessTokenStr(stringValue), nil
}

func (vo AccessTokenStr) String() string {
	return string(vo)
}
