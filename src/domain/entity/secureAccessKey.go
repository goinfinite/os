package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type SecureAccessKey struct {
	Id             valueObject.SecureAccessKeyId      `json:"id"`
	Name           valueObject.SecureAccessKeyName    `json:"name"`
	Content        valueObject.SecureAccessKeyContent `json:"-"`
	EncodedContent valueObject.EncodedContent         `json:"encodedContent"`
}

func NewSecureAccessKey(
	id valueObject.SecureAccessKeyId,
	name valueObject.SecureAccessKeyName,
	content valueObject.SecureAccessKeyContent,
	encodedContent valueObject.EncodedContent,
) SecureAccessKey {
	return SecureAccessKey{
		Id:             id,
		Name:           name,
		Content:        content,
		EncodedContent: encodedContent,
	}
}
