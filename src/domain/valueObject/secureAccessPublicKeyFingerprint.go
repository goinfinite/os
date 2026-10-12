package valueObject

import (
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var SecureAccessPublicKeyFingerprintRegex = regexp.MustCompile(`^SHA256:[\w\/\+\=]{43}$`)

type SecureAccessPublicKeyFingerprint string

func NewSecureAccessPublicKeyFingerprint(
	value any,
) (keyFingerprint SecureAccessPublicKeyFingerprint, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return keyFingerprint, errors.New("SecureAccessPublicKeyFingerprintMustBeString")
	}

	if !SecureAccessPublicKeyFingerprintRegex.MatchString(stringValue) {
		return keyFingerprint, errors.New("InvalidSecureAccessPublicKeyFingerprint")
	}

	return SecureAccessPublicKeyFingerprint(stringValue), nil
}

func (vo SecureAccessPublicKeyFingerprint) String() string {
	return string(vo)
}
