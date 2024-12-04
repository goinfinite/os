package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const SecureAccessPublicKeyContentRegex string = `^(?:ssh-(?:rsa|ed25519)) (?:[\w\/\+\=]+)(?: [\w@\-_]{6,32})?$`

type SecureAccessPublicKeyContent string

func NewSecureAccessPublicKeyContent(
	value interface{},
) (keyContent SecureAccessPublicKeyContent, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return keyContent, errors.New("SecureAccessPublicKeyContentMustBeString")
	}

	re := regexp.MustCompile(SecureAccessPublicKeyContentRegex)
	if !re.MatchString(stringValue) {
		return keyContent, errors.New("InvalidSecureAccessPublicKeyContent")
	}

	return SecureAccessPublicKeyContent(stringValue), nil
}

func (vo SecureAccessPublicKeyContent) String() string {
	return string(vo)
}

func (vo SecureAccessPublicKeyContent) ReadWithoutKeyName() string {
	keyContentParts := strings.Split(string(vo), " ")
	return keyContentParts[0] + " " + keyContentParts[1]
}

func (vo SecureAccessPublicKeyContent) ReadOnlyKeyName() (
	keyName SecureAccessPublicKeyName, err error,
) {
	keyContentParts := strings.Split(string(vo), " ")
	if len(keyContentParts) == 2 {
		return keyName, errors.New("SecureAccessPublicKeyContentHasNoName")
	}

	return NewSecureAccessPublicKeyName(keyContentParts[2])
}
