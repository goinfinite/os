package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdateCron struct {
	Id                valueObject.CronId         `json:"id"`
	Schedule          *valueObject.CronSchedule  `json:"schedule"`
	Command           *tkValueObject.UnixCommand `json:"command"`
	Comment           *valueObject.CronComment   `json:"comment"`
	ClearableFields   []string                   `json:"-"`
	OperatorAccountId tkValueObject.AccountId    `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress    `json:"-"`
}

func NewUpdateCron(
	id valueObject.CronId,
	schedule *valueObject.CronSchedule,
	command *tkValueObject.UnixCommand,
	comment *valueObject.CronComment,
	clearableFields []string,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
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
