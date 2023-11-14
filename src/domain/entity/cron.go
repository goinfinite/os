package entity

import "github.com/speedianet/os/src/domain/valueObject"

type Cron struct {
	Id       valueObject.CronId       `json:"id"`
	Schedule valueObject.CronSchedule `json:"schedule"`
	Command  valueObject.UnixCommand  `json:"command"`
	Comment  *valueObject.CronComment `json:"comment,omitempty"`
}

func NewCron(
	id valueObject.CronId,
	schedule valueObject.CronSchedule,
	command valueObject.UnixCommand,
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
		cronLineStr += " # " + cron.Comment.String() + "\n"
	}

	return cronLineStr
}
