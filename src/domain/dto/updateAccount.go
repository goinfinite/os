package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdateAccount struct {
	AccountId          *tkValueObject.AccountId `json:"accountId"`
	AccountUsername    *valueObject.Username    `json:"accountUsername"`
	Password           *tkValueObject.Password  `json:"password"`
	IsSuperAdmin       *bool                    `json:"isSuperAdmin"`
	ShouldUpdateApiKey *bool                    `json:"shouldUpdateApiKey"`
	OperatorAccountId  tkValueObject.AccountId  `json:"-"`
	OperatorIpAddress  tkValueObject.IpAddress  `json:"-"`
}

func NewUpdateAccount(
	accountId *tkValueObject.AccountId,
	accountUsername *valueObject.Username,
	password *tkValueObject.Password,
	isSuperAdmin, shouldUpdateApiKey *bool,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
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
