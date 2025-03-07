package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateAccount struct {
	AccountId          valueObject.AccountId `json:"accountId"`
	Password           *valueObject.Password `json:"password,omitempty"`
	IsSuperAdmin       *bool                 `json:"isSuperAdmin,omitempty"`
	ShouldUpdateApiKey *bool                 `json:"shouldUpdateApiKey,omitempty"`
	OperatorAccountId  valueObject.AccountId `json:"-"`
	OperatorIpAddress  valueObject.IpAddress `json:"-"`
}

func NewUpdateAccount(
	accountId valueObject.AccountId,
	password *valueObject.Password,
	isSuperAdmin *bool,
	shouldUpdateApiKey *bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateAccount {
	return UpdateAccount{
		AccountId:          accountId,
		Password:           password,
		IsSuperAdmin:       isSuperAdmin,
		ShouldUpdateApiKey: shouldUpdateApiKey,
		OperatorAccountId:  operatorAccountId,
		OperatorIpAddress:  operatorIpAddress,
	}
}
