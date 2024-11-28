package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateCron struct {
	Schedule          valueObject.CronSchedule `json:"schedule"`
	Command           valueObject.UnixCommand  `json:"command"`
	Comment           *valueObject.CronComment `json:"comment"`
	OperatorAccountId valueObject.AccountId    `json:"-"`
	OperatorIpAddress valueObject.IpAddress    `json:"-"`
}

func NewCreateCron(
	schedule valueObject.CronSchedule,
	command valueObject.UnixCommand,
	comment *valueObject.CronComment,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateCron {
	return CreateCron{
		Schedule:          schedule,
		Command:           command,
		Comment:           comment,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
