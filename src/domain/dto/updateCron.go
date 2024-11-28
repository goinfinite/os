package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateCron struct {
	Id                valueObject.CronId        `json:"id"`
	Schedule          *valueObject.CronSchedule `json:"schedule"`
	Command           *valueObject.UnixCommand  `json:"command"`
	Comment           *valueObject.CronComment  `json:"comment"`
	OperatorAccountId valueObject.AccountId     `json:"-"`
	OperatorIpAddress valueObject.IpAddress     `json:"-"`
}

func NewUpdateCron(
	id valueObject.CronId,
	schedule *valueObject.CronSchedule,
	command *valueObject.UnixCommand,
	comment *valueObject.CronComment,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateCron {
	return UpdateCron{
		Id:                id,
		Schedule:          schedule,
		Command:           command,
		Comment:           comment,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
