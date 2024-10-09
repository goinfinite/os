package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadScheduledTasksRequest struct {
	Pagination       Pagination                       `json:"pagination"`
	TaskId           *valueObject.ScheduledTaskId     `json:"taskId,omitempty"`
	TaskName         *valueObject.ScheduledTaskName   `json:"taskName,omitempty"`
	TaskStatus       *valueObject.ScheduledTaskStatus `json:"taskStatus,omitempty"`
	TaskTags         []valueObject.ScheduledTaskTag   `json:"taskTags,omitempty"`
	StartedBeforeAt  *valueObject.UnixTime            `json:"startedBeforeAt,omitempty"`
	StartedAfterAt   *valueObject.UnixTime            `json:"startedAfterAt,omitempty"`
	FinishedBeforeAt *valueObject.UnixTime            `json:"finishedBeforeAt,omitempty"`
	FinishedAfterAt  *valueObject.UnixTime            `json:"finishedAfterAt,omitempty"`
	CreatedBeforeAt  *valueObject.UnixTime            `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt   *valueObject.UnixTime            `json:"createdAfterAt,omitempty"`
}

type ReadScheduledTasksResponse struct {
	Pagination Pagination             `json:"pagination"`
	Tasks      []entity.ScheduledTask `json:"tasks"`
}
