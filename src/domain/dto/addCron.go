package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type AddCron struct {
	Schedule valueObject.CronSchedule `json:"schedule"`
	Command  valueObject.UnixCommand  `json:"command"`
	Comment  *valueObject.CronComment `json:"comment,omitempty"`
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
