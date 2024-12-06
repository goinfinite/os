package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
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
	fingerprint valueObject.SecureAccessPublicKeyFingerprint,
	name valueObject.SecureAccessPublicKeyName,
	createdAt, updatedAt valueObject.UnixTime,
) SecureAccessPublicKey {
	return SecureAccessPublicKey{
		Id:          id,
		AccountId:   accountId,
		Content:     content,
		Name:        name,
		Fingerprint: fingerprint,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
