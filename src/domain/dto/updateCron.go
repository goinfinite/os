package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type UpdateCron struct {
	Id       valueObject.CronId        `json:"id"`
	Schedule *valueObject.CronSchedule `json:"schedule,omitempty"`
	Command  *valueObject.UnixCommand  `json:"command,omitempty"`
	Comment  *valueObject.CronComment  `json:"comment,omitempty"`
}

func NewUpdateCron(
	id valueObject.CronId,
	schedule *valueObject.CronSchedule,
	command *valueObject.UnixCommand,
	comment *valueObject.CronComment,
) UpdateCron {
	return UpdateCron{
		Id:       id,
		Schedule: schedule,
		Command:  command,
		Comment:  comment,
	}
}
