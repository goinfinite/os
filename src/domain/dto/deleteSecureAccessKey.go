package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteSecureAccessKey struct {
	Id                valueObject.SecureAccessKeyId `json:"id"`
	OperatorAccountId valueObject.AccountId         `json:"-"`
	OperatorIpAddress valueObject.IpAddress         `json:"-"`
}

func NewDeleteSecureAccessKey(
	id valueObject.SecureAccessKeyId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteSecureAccessKey {
	return DeleteSecureAccessKey{
		Id:                id,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
