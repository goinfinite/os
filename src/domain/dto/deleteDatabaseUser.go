package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type DeleteDatabaseUser struct {
	DatabaseName      valueObject.DatabaseName    `json:"dbName"`
	Username          valueObject.DatabaseUsername `json:"username"`
	OperatorAccountId tkValueObject.AccountId     `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress     `json:"-"`
}

func NewDeleteDatabaseUser(
	dbName valueObject.DatabaseName,
	username valueObject.DatabaseUsername,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteDatabaseUser {
	return DeleteDatabaseUser{
		DatabaseName:      dbName,
		Username:          username,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
