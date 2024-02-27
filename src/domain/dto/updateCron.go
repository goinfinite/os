package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateCron struct {
	Id       valueObject.CronId        `json:"id"`
	Schedule *valueObject.CronSchedule `json:"schedule"`
	Command  *valueObject.UnixCommand  `json:"command"`
	Comment  *valueObject.CronComment  `json:"comment"`
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
