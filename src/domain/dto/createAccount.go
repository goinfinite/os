package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateAccount struct {
	Username          valueObject.Username    `json:"username"`
	Password          tkValueObject.Password  `json:"password"`
	IsSuperAdmin      bool                    `json:"isSuperAdmin"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

func NewCreateAccount(
	username valueObject.Username,
	password tkValueObject.Password,
	isSuperAdmin bool,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreateAccount {
	return CreateAccount{
		Username:          username,
		Password:          password,
		IsSuperAdmin:      isSuperAdmin,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
