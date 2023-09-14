package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type UpdateCron struct {
	Schedule valueObject.CronSchedule `json:"schedule"`
	Command  *valueObject.UnixCommand `json:"command,omitempty"`
	Comment  *valueObject.CronComment `json:"comment,omitempty"`
}

func NewUpdateCron(
	schedule valueObject.CronSchedule,
	command *valueObject.UnixCommand,
	comment *valueObject.CronComment,
) UpdateCron {
	return UpdateCron{
		Schedule: schedule,
		Command:  command,
		Comment:  comment,
	}
}
