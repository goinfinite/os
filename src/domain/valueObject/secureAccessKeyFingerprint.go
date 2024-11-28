package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const secureAccessKeyFingerprintRegex string = `^SHA256:[\w\/\+\=]{43}$`

type SecureAccessKeyFingerprint string

func NewSecureAccessKeyFingerprint(
	value interface{},
) (keyFingerprint SecureAccessKeyFingerprint, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return keyFingerprint, errors.New("SecureAccessKeyFingerprintMustBeString")
	}

	re := regexp.MustCompile(secureAccessKeyFingerprintRegex)
	if !re.MatchString(stringValue) {
		return keyFingerprint, errors.New("InvalidSecureAccessKeyFingerprint")
	}

	return SecureAccessKeyFingerprint(stringValue), nil
}

func (vo SecureAccessKeyFingerprint) String() string {
	return string(vo)
}
