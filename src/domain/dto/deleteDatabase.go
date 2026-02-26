package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type DeleteDatabase struct {
	DatabaseName      valueObject.DatabaseName `json:"dbName"`
	OperatorAccountId tkValueObject.AccountId  `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress  `json:"-"`
}

func NewDeleteDatabase(
	dbName valueObject.DatabaseName,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteDatabase {
	return DeleteDatabase{
		DatabaseName:      dbName,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
