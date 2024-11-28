package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteDatabase struct {
	DatabaseName      valueObject.DatabaseName `json:"dbName"`
	OperatorAccountId valueObject.AccountId    `json:"-"`
	OperatorIpAddress valueObject.IpAddress    `json:"-"`
}

func NewDeleteDatabase(
	dbName valueObject.DatabaseName,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteDatabase {
	return DeleteDatabase{
		DatabaseName:      dbName,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
