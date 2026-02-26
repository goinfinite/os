package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type DeleteSecureAccessPublicKey struct {
	Id                valueObject.SecureAccessPublicKeyId `json:"id"`
	OperatorAccountId tkValueObject.AccountId             `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress             `json:"-"`
}

func NewDeleteSecureAccessPublicKey(
	id valueObject.SecureAccessPublicKeyId,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteSecureAccessPublicKey {
	return DeleteSecureAccessPublicKey{
		Id:                id,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
