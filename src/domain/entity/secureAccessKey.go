package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type SecureAccessKey struct {
	Id        valueObject.SecureAccessKeyId      `json:"id"`
	AccountId valueObject.AccountId              `json:"accountId"`
	Name      valueObject.SecureAccessKeyName    `json:"name"`
	Content   valueObject.SecureAccessKeyContent `json:"content"`
	CreatedAt valueObject.UnixTime               `json:"createdAt"`
	UpdatedAt valueObject.UnixTime               `json:"updatedAt"`
}

func NewSecureAccessKey(
	id valueObject.SecureAccessKeyId,
	accountId valueObject.AccountId,
	name valueObject.SecureAccessKeyName,
	content valueObject.SecureAccessKeyContent,
	createdAt, updatedAt valueObject.UnixTime,
) SecureAccessKey {
	return SecureAccessKey{
		Id:        id,
		AccountId: accountId,
		Name:      name,
		Content:   content,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
