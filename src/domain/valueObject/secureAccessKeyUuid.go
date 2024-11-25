package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const secureAccessKeyUuidRegex string = `^\w{10,16}$`

type SecureAccessKeyUuid string

func NewSecureAccessKeyUuid(
	value interface{},
) (keyUuid SecureAccessKeyUuid, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return keyUuid, errors.New("SecureAccessKeyUuidMustBeString")
	}

	re := regexp.MustCompile(secureAccessKeyUuidRegex)
	if !re.MatchString(stringValue) {
		return keyUuid, errors.New("InvalidSecureAccessKeyUuid")
	}

	return SecureAccessKeyUuid(stringValue), nil
}

func (vo SecureAccessKeyUuid) String() string {
	return string(vo)
}
