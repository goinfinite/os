package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateDatabaseUser struct {
	DatabaseName      valueObject.DatabaseName        `json:"dbName"`
	Username          valueObject.DatabaseUsername    `json:"username"`
	Password          valueObject.Password            `json:"password"`
	Privileges        []valueObject.DatabasePrivilege `json:"privileges"`
	OperatorAccountId valueObject.AccountId           `json:"-"`
	OperatorIpAddress valueObject.IpAddress           `json:"-"`
}

func NewCreateDatabaseUser(
	dbName valueObject.DatabaseName,
	username valueObject.DatabaseUsername,
	password valueObject.Password,
	privileges []valueObject.DatabasePrivilege,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
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
