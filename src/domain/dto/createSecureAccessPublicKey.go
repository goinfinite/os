package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateSecureAccessPublicKey struct {
	AccountId         valueObject.AccountId                    `json:"accountId"`
	Content           valueObject.SecureAccessPublicKeyContent `json:"content"`
	Name              valueObject.SecureAccessPublicKeyName    `json:"name"`
	OperatorAccountId valueObject.AccountId                    `json:"-"`
	OperatorIpAddress valueObject.IpAddress                    `json:"-"`
}

func NewCreateSecureAccessPublicKey(
	accountId valueObject.AccountId,
	content valueObject.SecureAccessPublicKeyContent,
	name valueObject.SecureAccessPublicKeyName,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateSecureAccessPublicKey {
	return CreateSecureAccessPublicKey{
		AccountId:         accountId,
		Content:           content,
		Name:              name,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
