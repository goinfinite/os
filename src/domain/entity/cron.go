package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type Cron struct {
	Id       valueObject.CronId        `json:"id"`
	Schedule valueObject.CronSchedule  `json:"schedule"`
	Command  tkValueObject.UnixCommand `json:"command"`
	Comment  *valueObject.CronComment  `json:"comment"`
}

func NewCron(
	id valueObject.CronId,
	schedule valueObject.CronSchedule,
	command tkValueObject.UnixCommand,
	comment *valueObject.CronComment,
) Cron {
	return Cron{
		Id:       id,
		Schedule: schedule,
		Command:  command,
		Comment:  comment,
	}
}

func (cron Cron) String() string {
	cronLineStr := cron.Schedule.String() + " " + cron.Command.String()
	if cron.Comment != nil {
		cronLineStr += " # " + cron.Comment.String()
	}

	return cronLineStr
}
