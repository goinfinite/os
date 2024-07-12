package dto

import "github.com/speedianet/os/src/domain/valueObject"

type DeleteCron struct {
	Id      *valueObject.CronId      `json:"id"`
	Comment *valueObject.CronComment `json:"comment"`
}

func NewDeleteCron(
	id *valueObject.CronId,
	comment *valueObject.CronComment,
) DeleteCron {
	return DeleteCron{
		Id:      id,
		Comment: comment,
	}
}
