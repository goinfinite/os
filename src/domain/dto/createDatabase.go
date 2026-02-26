package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateDatabase struct {
	DatabaseName      valueObject.DatabaseName `json:"dbName"`
	OperatorAccountId tkValueObject.AccountId  `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress  `json:"-"`
}

func NewCreateDatabase(
	dbName valueObject.DatabaseName,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreateDatabase {
	return CreateDatabase{
		DatabaseName:      dbName,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
