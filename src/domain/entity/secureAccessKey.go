package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type SecureAccessKey struct {
	Name    valueObject.SecureAccessKeyName    `json:"secureAccessKeyName"`
	Content valueObject.SecureAccessKeyContent `json:"secureAccessKeyContent"`
}

func NewSecureAccessKey(
	name valueObject.SecureAccessKeyName,
	content valueObject.SecureAccessKeyContent,
) SecureAccessKey {
	return SecureAccessKey{
		Name:    name,
		Content: content,
	}
}
