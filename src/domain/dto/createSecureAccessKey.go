package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateSecureAccessKey struct {
	Name              valueObject.SecureAccessKeyName    `json:"name"`
	Content           valueObject.SecureAccessKeyContent `json:"content"`
	AccountId         valueObject.AccountId              `json:"accountId"`
	OperatorAccountId valueObject.AccountId              `json:"-"`
	OperatorIpAddress valueObject.IpAddress              `json:"-"`
}

func NewCreateSecureAccessKey(
	name valueObject.SecureAccessKeyName,
	content valueObject.SecureAccessKeyContent,
	accountId, operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateSecureAccessKey {
	return CreateSecureAccessKey{
		Name:              name,
		Content:           content,
		AccountId:         accountId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
