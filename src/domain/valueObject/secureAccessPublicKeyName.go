package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const SecureAccessPublicKeyNameRegex string = `^[\w@\-_]{6,32}$`

type SecureAccessPublicKeyName string

func NewSecureAccessPublicKeyName(
	value interface{},
) (keyName SecureAccessPublicKeyName, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return keyName, errors.New("SecureAccessPublicKeyNameMustBeString")
	}

	re := regexp.MustCompile(SecureAccessPublicKeyNameRegex)
	if !re.MatchString(stringValue) {
		return keyName, errors.New("InvalidSecureAccessPublicKeyName")
	}

	return SecureAccessPublicKeyName(stringValue), nil
}

func (vo SecureAccessPublicKeyName) String() string {
	return string(vo)
}
