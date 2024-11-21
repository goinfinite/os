package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const secureAccessKeyContentRegex string = `^(?:ssh-rsa) (?:[\w\/\+\=]+)(?: [\w@\-_]{6,32})?$`

type SecureAccessKeyContent string

func NewSecureAccessKeyContent(
	value interface{},
) (keyContent SecureAccessKeyContent, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return keyContent, errors.New("SecureAccessKeyContentMustBeString")
	}

	re := regexp.MustCompile(secureAccessKeyContentRegex)
	if !re.MatchString(stringValue) {
		return keyContent, errors.New("InvalidSecureAccessKeyContent")
	}

	return SecureAccessKeyContent(stringValue), nil
}

func (vo SecureAccessKeyContent) String() string {
	return string(vo)
}

func (vo SecureAccessKeyContent) ReadWithoutKeyName() string {
	keyContentParts := strings.Split(string(vo), " ")
	return keyContentParts[0] + " " + keyContentParts[1]
}

func (vo SecureAccessKeyContent) ReadOnlyKeyName() (
	keyName SecureAccessKeyName, err error,
) {
	keyContentParts := strings.Split(string(vo), " ")
	if len(keyContentParts) == 2 {
		return keyName, errors.New("SecureAccessKeyNameNotFound")
	}

	return NewSecureAccessKeyName(keyContentParts[2])
}
