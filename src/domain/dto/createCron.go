package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateCron struct {
	Schedule          valueObject.CronSchedule `json:"schedule"`
	Command           tkValueObject.UnixCommand  `json:"command"`
	Comment           *valueObject.CronComment `json:"comment"`
	OperatorAccountId tkValueObject.AccountId    `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress    `json:"-"`
}

func NewCreateCron(
	schedule valueObject.CronSchedule,
	command tkValueObject.UnixCommand,
	comment *valueObject.CronComment,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreateCron {
	return CreateCron{
		Schedule:          schedule,
		Command:           command,
		Comment:           comment,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
