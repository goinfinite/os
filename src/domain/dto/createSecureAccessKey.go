package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateSecureAccessKey struct {
	AccountId         valueObject.AccountId              `json:"accountId"`
	Content           valueObject.SecureAccessKeyContent `json:"content"`
	Name              valueObject.SecureAccessKeyName    `json:"name"`
	OperatorAccountId valueObject.AccountId              `json:"-"`
	OperatorIpAddress valueObject.IpAddress              `json:"-"`
}

func NewCreateSecureAccessKey(
	accountId valueObject.AccountId,
	content valueObject.SecureAccessKeyContent,
	name valueObject.SecureAccessKeyName,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateSecureAccessKey {
	return CreateSecureAccessKey{
		AccountId:         accountId,
		Content:           content,
		Name:              name,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
