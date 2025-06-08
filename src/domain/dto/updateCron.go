package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateCron struct {
	Id                valueObject.CronId        `json:"id"`
	Schedule          *valueObject.CronSchedule `json:"schedule"`
	Command           *valueObject.UnixCommand  `json:"command"`
	Comment           *valueObject.CronComment  `json:"comment"`
	ClearableFields   []string                  `json:"-"`
	OperatorAccountId valueObject.AccountId     `json:"-"`
	OperatorIpAddress valueObject.IpAddress     `json:"-"`
}

func NewUpdateCron(
	id valueObject.CronId,
	schedule *valueObject.CronSchedule,
	command *valueObject.UnixCommand,
	comment *valueObject.CronComment,
	clearableFields []string,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateCron {
	return UpdateCron{
		Id:                id,
		Schedule:          schedule,
		Command:           command,
		Comment:           comment,
		ClearableFields:   clearableFields,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
