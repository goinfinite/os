package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type SecureAccessPublicKey struct {
	Id          valueObject.SecureAccessPublicKeyId          `json:"id"`
	AccountId   valueObject.AccountId                        `json:"accountId"`
	Name        valueObject.SecureAccessPublicKeyName        `json:"name"`
	Content     valueObject.SecureAccessPublicKeyContent     `json:"-"`
	Fingerprint valueObject.SecureAccessPublicKeyFingerprint `json:"fingerprint"`
	CreatedAt   valueObject.UnixTime                         `json:"createdAt"`
	UpdatedAt   valueObject.UnixTime                         `json:"updatedAt"`
}

func NewSecureAccessPublicKey(
	id valueObject.SecureAccessPublicKeyId,
	accountId valueObject.AccountId,
	name valueObject.SecureAccessPublicKeyName,
	content valueObject.SecureAccessPublicKeyContent,
	fingerprint valueObject.SecureAccessPublicKeyFingerprint,
	createdAt, updatedAt valueObject.UnixTime,
) SecureAccessPublicKey {
	return SecureAccessPublicKey{
		Id:          id,
		AccountId:   accountId,
		Name:        name,
		Content:     content,
		Fingerprint: fingerprint,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
