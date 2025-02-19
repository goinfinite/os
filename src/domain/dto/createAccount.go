package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateAccount struct {
	Username          valueObject.Username  `json:"username"`
	Password          valueObject.Password  `json:"password"`
	IsSuperAdmin      bool                  `json:"isSupermanAdmin"`
	OperatorAccountId valueObject.AccountId `json:"-"`
	OperatorIpAddress valueObject.IpAddress `json:"-"`
}

func NewCreateAccount(
	username valueObject.Username,
	password valueObject.Password,
	isSuperAdmin bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateAccount {
	return CreateAccount{
		Username:          username,
		Password:          password,
		IsSuperAdmin:      isSuperAdmin,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
