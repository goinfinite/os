package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type SecureAccessKey struct {
	HashId  valueObject.SecureAccessKeyHashId  `json:"hashId"`
	Name    valueObject.SecureAccessKeyName    `json:"name"`
	Content valueObject.SecureAccessKeyContent `json:"-"`
}

func NewSecureAccessKey(
	hashId valueObject.SecureAccessKeyHashId,
	name valueObject.SecureAccessKeyName,
	content valueObject.SecureAccessKeyContent,
) SecureAccessKey {
	return SecureAccessKey{
		HashId:  hashId,
		Name:    name,
		Content: content,
	}
}
