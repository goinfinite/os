package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteSecureAccessKey struct {
	Id                valueObject.SecureAccessKeyId `json:"id"`
	AccountId         valueObject.AccountId         `json:"-"`
	OperatorAccountId valueObject.AccountId         `json:"-"`
	OperatorIpAddress valueObject.IpAddress         `json:"-"`
}

func NewDeleteSecureAccessKey(
	id valueObject.SecureAccessKeyId,
	accountId, operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteSecureAccessKey {
	return DeleteSecureAccessKey{
		Id:                id,
		AccountId:         accountId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
