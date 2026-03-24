package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateSecureAccessPublicKey struct {
	AccountId         tkValueObject.AccountId                  `json:"accountId"`
	Content           valueObject.SecureAccessPublicKeyContent `json:"content"`
	Name              valueObject.SecureAccessPublicKeyName    `json:"name"`
	OperatorAccountId tkValueObject.AccountId                  `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress                  `json:"-"`
}

func NewCreateSecureAccessPublicKey(
	accountId tkValueObject.AccountId,
	content valueObject.SecureAccessPublicKeyContent,
	name valueObject.SecureAccessPublicKeyName,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreateSecureAccessPublicKey {
	return CreateSecureAccessPublicKey{
		AccountId:         accountId,
		Content:           content,
		Name:              name,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
