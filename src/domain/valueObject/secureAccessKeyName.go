package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const secureAccessKeyNameRegex string = `^[\w@\-_]{6,32}$`

type SecureAccessKeyName string

func NewSecureAccessKeyName(
	value interface{},
) (keyName SecureAccessKeyName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return keyName, errors.New("SecureAccessKeyNameMustBeString")
	}

	re := regexp.MustCompile(secureAccessKeyNameRegex)
	if !re.MatchString(stringValue) {
		return keyName, errors.New("InvalidSecureAccessKeyName")
	}

	return SecureAccessKeyName(stringValue), nil
}

func (vo SecureAccessKeyName) String() string {
	return string(vo)
}
