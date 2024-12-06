package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	"golang.org/x/crypto/ssh"
)

type SecureAccessPublicKeyContent struct {
	Content     string
	Fingerprint string
}

func NewSecureAccessPublicKeyContent(
	value interface{},
) (keyContent SecureAccessPublicKeyContent, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return keyContent, errors.New("SecureAccessPublicKeyContentMustBeString")
	}

	publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(stringValue))
	if err != nil {
		return keyContent, errors.New("InvalidSecureAccessPublicKeyContent")
	}

	return SecureAccessPublicKeyContent{
		Content:     stringValue,
		Fingerprint: ssh.FingerprintSHA256(publicKey),
	}, nil
}

func (vo SecureAccessPublicKeyContent) String() string {
	return string(vo.Content)
}

func (vo SecureAccessPublicKeyContent) ReadWithoutKeyName() string {
	keyContentParts := strings.Split(string(vo.Content), " ")
	return keyContentParts[0] + " " + keyContentParts[1]
}

func (vo SecureAccessPublicKeyContent) ReadOnlyKeyName() (
	keyName SecureAccessPublicKeyName, err error,
) {
	keyContentParts := strings.Split(string(vo.Content), " ")
	if len(keyContentParts) == 2 {
		return keyName, errors.New("SecureAccessPublicKeyContentHasNoName")
	}

	return NewSecureAccessPublicKeyName(keyContentParts[2])
}

func (vo SecureAccessPublicKeyContent) ReadFingerprint() (
	SecureAccessPublicKeyFingerprint, error,
) {
	return NewSecureAccessPublicKeyFingerprint(vo.Fingerprint)
}
