package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteCron struct {
	Id                *valueObject.CronId      `json:"id"`
	Comment           *valueObject.CronComment `json:"comment"`
	OperatorAccountId valueObject.AccountId    `json:"-"`
	OperatorIpAddress valueObject.IpAddress    `json:"-"`
}

func NewDeleteCron(
	id *valueObject.CronId,
	comment *valueObject.CronComment,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteCron {
	return DeleteCron{
		Id:                id,
		Comment:           comment,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
