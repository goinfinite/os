package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var SecureAccessPublicKeyNameRegex = regexp.MustCompile(`^[A-Za-z0-9][\w@\-_]{5,32}$`)

type SecureAccessPublicKeyName string

func NewSecureAccessPublicKeyName(
	value any,
) (keyName SecureAccessPublicKeyName, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return keyName, errors.New("SecureAccessPublicKeyNameMustBeString")
	}

	if !SecureAccessPublicKeyNameRegex.MatchString(stringValue) {
		return keyName, errors.New("InvalidSecureAccessPublicKeyName")
	}

	return SecureAccessPublicKeyName(stringValue), nil
}

func (vo SecureAccessPublicKeyName) String() string {
	return string(vo)
}
