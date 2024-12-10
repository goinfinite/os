package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const SecureAccessPublicKeyFingerprintRegex string = `^SHA256:[\w\/\+\=]{43}$`

type SecureAccessPublicKeyFingerprint string

func NewSecureAccessPublicKeyFingerprint(
	value interface{},
) (keyFingerprint SecureAccessPublicKeyFingerprint, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return keyFingerprint, errors.New("SecureAccessPublicKeyFingerprintMustBeString")
	}

	re := regexp.MustCompile(SecureAccessPublicKeyFingerprintRegex)
	if !re.MatchString(stringValue) {
		return keyFingerprint, errors.New("InvalidSecureAccessPublicKeyFingerprint")
	}

	return SecureAccessPublicKeyFingerprint(stringValue), nil
}

func (vo SecureAccessPublicKeyFingerprint) String() string {
	return string(vo)
}
