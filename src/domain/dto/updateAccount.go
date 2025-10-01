package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateAccount struct {
	AccountId          *valueObject.AccountId `json:"accountId"`
	AccountUsername    *valueObject.Username  `json:"accountUsername"`
	Password           *valueObject.Password  `json:"password"`
	IsSuperAdmin       *bool                  `json:"isSuperAdmin"`
	ShouldUpdateApiKey *bool                  `json:"shouldUpdateApiKey"`
	OperatorAccountId  valueObject.AccountId  `json:"-"`
	OperatorIpAddress  valueObject.IpAddress  `json:"-"`
}

func NewUpdateAccount(
	accountId *valueObject.AccountId,
	accountUsername *valueObject.Username,
	password *valueObject.Password,
	isSuperAdmin, shouldUpdateApiKey *bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateAccount {
	return UpdateAccount{
		AccountId:          accountId,
		AccountUsername:    accountUsername,
		Password:           password,
		IsSuperAdmin:       isSuperAdmin,
		ShouldUpdateApiKey: shouldUpdateApiKey,
		OperatorAccountId:  operatorAccountId,
		OperatorIpAddress:  operatorIpAddress,
	}
}
