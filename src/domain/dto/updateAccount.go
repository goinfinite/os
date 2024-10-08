package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateAccount struct {
	AccountId          valueObject.AccountId `json:"accountId"`
	Password           *valueObject.Password `json:"password,omitempty"`
	ShouldUpdateApiKey *bool                 `json:"shouldUpdateApiKey,omitempty"`
	OperatorAccountId  valueObject.AccountId `json:"-"`
	OperatorIpAddress  valueObject.IpAddress `json:"-"`
}

func NewUpdateAccount(
	accountId valueObject.AccountId,
	password *valueObject.Password,
	shouldUpdateApiKey *bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateAccount {
	return UpdateAccount{
		AccountId:          accountId,
		Password:           password,
		ShouldUpdateApiKey: shouldUpdateApiKey,
		OperatorAccountId:  operatorAccountId,
		OperatorIpAddress:  operatorIpAddress,
	}
}
