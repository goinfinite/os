package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteDatabaseUser struct {
	DatabaseName      valueObject.DatabaseName     `json:"dbName"`
	Username          valueObject.DatabaseUsername `json:"username"`
	OperatorAccountId valueObject.AccountId        `json:"-"`
	OperatorIpAddress valueObject.IpAddress        `json:"-"`
}

func NewDeleteDatabaseUser(
	dbName valueObject.DatabaseName,
	username valueObject.DatabaseUsername,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteDatabaseUser {
	return DeleteDatabaseUser{
		DatabaseName:      dbName,
		Username:          username,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
