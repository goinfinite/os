package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateDatabaseUser struct {
	DatabaseName      valueObject.DatabaseName        `json:"dbName"`
	Username          valueObject.DatabaseUsername     `json:"username"`
	Password          tkValueObject.WeakPassword      `json:"password"`
	Privileges        []valueObject.DatabasePrivilege `json:"privileges"`
	OperatorAccountId tkValueObject.AccountId         `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress         `json:"-"`
}

func NewCreateDatabaseUser(
	dbName valueObject.DatabaseName,
	username valueObject.DatabaseUsername,
	password tkValueObject.WeakPassword,
	privileges []valueObject.DatabasePrivilege,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreateDatabaseUser {
	return CreateDatabaseUser{
		DatabaseName:      dbName,
		Username:          username,
		Password:          password,
		Privileges:        privileges,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
