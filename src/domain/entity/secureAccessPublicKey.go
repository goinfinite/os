package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type SecureAccessPublicKey struct {
	Id          valueObject.SecureAccessPublicKeyId          `json:"id"`
	AccountId   tkValueObject.AccountId                      `json:"accountId"`
	Content     valueObject.SecureAccessPublicKeyContent     `json:"-"`
	Name        valueObject.SecureAccessPublicKeyName        `json:"name"`
	Fingerprint valueObject.SecureAccessPublicKeyFingerprint `json:"fingerprint"`
	CreatedAt   tkValueObject.UnixTime                       `json:"createdAt"`
	UpdatedAt   tkValueObject.UnixTime                       `json:"updatedAt"`
}

func NewSecureAccessPublicKey(
	id valueObject.SecureAccessPublicKeyId,
	accountId tkValueObject.AccountId,
	content valueObject.SecureAccessPublicKeyContent,
	fingerprint valueObject.SecureAccessPublicKeyFingerprint,
	name valueObject.SecureAccessPublicKeyName,
	createdAt, updatedAt tkValueObject.UnixTime,
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
