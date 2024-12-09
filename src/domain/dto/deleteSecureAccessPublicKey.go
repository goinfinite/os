package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteSecureAccessPublicKey struct {
	Id                valueObject.SecureAccessPublicKeyId `json:"id"`
	OperatorAccountId valueObject.AccountId               `json:"-"`
	OperatorIpAddress valueObject.IpAddress               `json:"-"`
}

func NewDeleteSecureAccessPublicKey(
	id valueObject.SecureAccessPublicKeyId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteSecureAccessPublicKey {
	return DeleteSecureAccessPublicKey{
		Id:                id,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
