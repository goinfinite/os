package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateAccount struct {
	Username          valueObject.Username  `json:"username"`
	Password          valueObject.Password  `json:"password"`
	OperatorAccountId valueObject.AccountId `json:"-"`
	OperatorIpAddress valueObject.IpAddress `json:"-"`
}

func NewCreateAccount(
	username valueObject.Username,
	password valueObject.Password,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateAccount {
	return CreateAccount{
		Username:          username,
		Password:          password,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
