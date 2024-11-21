package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const secureAccessKeyHashIdRegex string = `^[\w\/\+\=]{24,}$`

type SecureAccessKeyHashId string

func NewSecureAccessKeyHashId(
	value interface{},
) (keyHashId SecureAccessKeyHashId, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return keyHashId, errors.New("SecureAccessKeyHashIdMustBeString")
	}

	re := regexp.MustCompile(secureAccessKeyHashIdRegex)
	if !re.MatchString(stringValue) {
		return keyHashId, errors.New("InvalidSecureAccessKeyHashId")
	}

	return SecureAccessKeyHashId(stringValue[:24]), nil
}

func (vo SecureAccessKeyHashId) String() string {
	return string(vo)
}
