package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateDatabase struct {
	DatabaseName      valueObject.DatabaseName `json:"dbName"`
	DatabaseType      valueObject.DatabaseType `json:"dbType"`
	OperatorAccountId valueObject.AccountId    `json:"-"`
	OperatorIpAddress valueObject.IpAddress    `json:"-"`
}

func NewCreateDatabase(
	dbName valueObject.DatabaseName,
	dbType valueObject.DatabaseType,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateDatabase {
	return CreateDatabase{
		DatabaseName:      dbName,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
