package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type DeleteMapping struct {
	MappingId         valueObject.MappingId   `json:"mappingId"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

func NewDeleteMapping(
	mappingId valueObject.MappingId,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteMapping {
	return DeleteMapping{
		MappingId:         mappingId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
