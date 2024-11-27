package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteMapping struct {
	MappingId         valueObject.MappingId `json:"mappingId"`
	OperatorAccountId valueObject.AccountId `json:"-"`
	OperatorIpAddress valueObject.IpAddress `json:"-"`
}

func NewDeleteMapping(
	mappingId valueObject.MappingId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteMapping {
	return DeleteMapping{
		MappingId:         mappingId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
