package dto

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type DeleteAccount struct {
	AccountId         tkValueObject.AccountId `json:"accountId"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

func NewDeleteAccount(
	accountId, operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteAccount {
	return DeleteAccount{
		AccountId:         accountId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
