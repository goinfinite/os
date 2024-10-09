package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateCron struct {
	Schedule valueObject.CronSchedule `json:"schedule"`
	Command  valueObject.UnixCommand  `json:"command"`
	Comment  *valueObject.CronComment `json:"comment"`
}

func NewCreateCron(
	schedule valueObject.CronSchedule,
	command valueObject.UnixCommand,
	comment *valueObject.CronComment,
) CreateCron {
	return CreateCron{
		Schedule: schedule,
		Command:  command,
		Comment:  comment,
	}
}
