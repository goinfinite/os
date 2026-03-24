package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ReadScheduledTasksRequest struct {
	Pagination       tkDto.Pagination                 `json:"pagination"`
	TaskId           *valueObject.ScheduledTaskId     `json:"taskId,omitempty"`
	TaskName         *valueObject.ScheduledTaskName   `json:"taskName,omitempty"`
	TaskStatus       *valueObject.ScheduledTaskStatus `json:"taskStatus,omitempty"`
	TaskTags         []valueObject.ScheduledTaskTag   `json:"taskTags,omitempty"`
	StartedBeforeAt  *tkValueObject.UnixTime          `json:"startedBeforeAt,omitempty"`
	StartedAfterAt   *tkValueObject.UnixTime          `json:"startedAfterAt,omitempty"`
	FinishedBeforeAt *tkValueObject.UnixTime          `json:"finishedBeforeAt,omitempty"`
	FinishedAfterAt  *tkValueObject.UnixTime          `json:"finishedAfterAt,omitempty"`
	CreatedBeforeAt  *tkValueObject.UnixTime          `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt   *tkValueObject.UnixTime          `json:"createdAfterAt,omitempty"`
}

type ReadScheduledTasksResponse struct {
	Pagination tkDto.Pagination       `json:"pagination"`
	Tasks      []entity.ScheduledTask `json:"tasks"`
}
