package entity

import (
	"errors"

	"github.com/goinfinite/os/src/domain/valueObject"
	"golang.org/x/crypto/ssh"
)

type SecureAccessPublicKey struct {
	Id          valueObject.SecureAccessPublicKeyId          `json:"id"`
	AccountId   valueObject.AccountId                        `json:"accountId"`
	Content     valueObject.SecureAccessPublicKeyContent     `json:"-"`
	Name        valueObject.SecureAccessPublicKeyName        `json:"name"`
	Fingerprint valueObject.SecureAccessPublicKeyFingerprint `json:"fingerprint"`
	CreatedAt   valueObject.UnixTime                         `json:"createdAt"`
	UpdatedAt   valueObject.UnixTime                         `json:"updatedAt"`
}

func NewSecureAccessPublicKey(
	id valueObject.SecureAccessPublicKeyId,
	accountId valueObject.AccountId,
	content valueObject.SecureAccessPublicKeyContent,
	namePtr *valueObject.SecureAccessPublicKeyName,
	createdAt, updatedAt valueObject.UnixTime,
) (secureAccessPublicKey SecureAccessPublicKey, err error) {
	contentBytes := []byte(content.String())
	publicKey, publicKeyNameStr, _, _, err := ssh.ParseAuthorizedKey(contentBytes)
	if err != nil {
		return secureAccessPublicKey, errors.New("SecureAccessPublicKeyParseError")
	}

	if namePtr != nil {
		publicKeyNameStr = namePtr.String()
	}
	name, err := valueObject.NewSecureAccessPublicKeyName(publicKeyNameStr)
	if err != nil {
		return secureAccessPublicKey, err
	}

	fingerprintStr := ssh.FingerprintSHA256(publicKey)
	fingerprint, err := valueObject.NewSecureAccessPublicKeyFingerprint(fingerprintStr)
	if err != nil {
		return secureAccessPublicKey, err
	}

	return SecureAccessPublicKey{
		Id:          id,
		AccountId:   accountId,
		Name:        name,
		Content:     content,
		Fingerprint: fingerprint,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
