package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type DeleteCron struct {
	Id                *valueObject.CronId      `json:"id"`
	Comment           *valueObject.CronComment `json:"comment"`
	OperatorAccountId tkValueObject.AccountId  `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress  `json:"-"`
}

func NewDeleteCron(
	id *valueObject.CronId,
	comment *valueObject.CronComment,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteCron {
	return DeleteCron{
		Id:                id,
		Comment:           comment,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
