package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddCron struct {
	Schedule valueObject.CronSchedule `json:"schedule"`
	Command  valueObject.UnixCommand  `json:"command"`
	Comment  *valueObject.CronComment `json:"comment"`
}

func NewAddCron(
	schedule valueObject.CronSchedule,
	command valueObject.UnixCommand,
	comment *valueObject.CronComment,
) AddCron {
	return AddCron{
		Schedule: schedule,
		Command:  command,
		Comment:  comment,
	}
}
